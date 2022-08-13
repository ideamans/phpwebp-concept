<?php

function parse_request(&$method, &$image_file_path, &$cache_file_path, &$cache_key) {
  // リクエストメソッド
  $method = $_SERVER['REQUEST_METHOD'];

  // 画像ファイルの実パスを取得
  $document_root_path = $_SERVER['DOCUMENT_ROOT'];
  $image_requested_path = $_SERVER['REDIRECT_URL'];
  $image_file_path = realpath(rtrim($document_root_path, '/') . $image_requested_path);

  // 実パスがなければ404を返す
  if (!$image_file_path) throw new \Exception('404 Not Found');

  // 実パスがDOCUMENT_ROOT以下でなければ403を返す
  // これは悪意のあるDOCUMENT_ROOT以外のファイルへのリクエストを防ぐため
  if (strpos($image_file_path, $document_root_path) !== 0) throw new \Exception('403 Forbidden');

  // キャッシュキーとファイルパスを導出: sha1(リクエストパス\t更新日\tファイルサイズ)
  // これは元画像の変更(更新日またはファイルサイズの変更)を検知し速やか再変換するため
  $cache_dir_path = rtrim(sys_get_temp_dir(), '/') . '/phpwebp-concept';
  // mkdirは他のリクエストとタイミングが重なることがあるので@でエラーを抑制する
  @mkdir($cache_dir_path, 0777, true);

  $cache_seed = [$image_file_path, filemtime($image_file_path), filesize($image_file_path)];
  $cache_key = sha1(implode("\t", $cache_seed));
  $cache_file_path = rtrim($cache_dir_path, '/') . '/' . $cache_key;
}

/**
 * WebP変換のためのコマンド実行を安全に行う
 */
function safe_webp_exec_via_stdin($command, $args, $input_file_path, &$stdout = '', &$stderr = '') {
  // コマンドパスを特定する(bin/$os-$arch)
  // linux-x86_64は標準で同梱するが、他のOS/アーキテクチャについては各自、cwebpとgif2webpを設置する
  $arch = strtolower(PHP_OS . '-' . php_uname('m'));
  $ext = substr(PHP_OS, 0, 3) === 'WIN' ? '.ext' : '';
  $command_path = __DIR__ . "/bin/$arch/$command$ext";
  $shell = "$command_path $args";

  // $input_file_path を標準入力から渡してWebPへの変換を行う
  // これは$input_file_pathが、改変の可能性は低いもののユーザー由来の値であり、
  // コマンド引数に渡すことがコマンドインジェクションにつながる恐れがあるため
  $desc = [
    fopen($input_file_path, 'rb'),
    ['pipe', 'w'],
    ['pipe', 'w'],
  ];
  $cwd = getcwd();
  $env = null;
  $code = -1;
  $proc = proc_open($shell, $desc, $pipes, $cwd, $env);

  if (is_resource($proc)) {
    $stdout = stream_get_contents($pipes[1]);
    $stderr = stream_get_contents($pipes[2]);
    $code = proc_close($proc);
  } else {
    throw new \Exception("Failed to open proc: $shell");
  }

  if ($code !== 0) {
    throw new \Exception("Command failed: $shell < $input_file_path; Exit code: $code; Stdout: $stdout; Stderr: $stderr;");
  }
}

function handle_as_http_exception($ex) {
  $message = $ex->getMessage();

  // エラーメッセージが数字3桁から始まる場合はHTTPステータスとして返す
  if (preg_match('/^\d{3}/', $message)) {
    header("HTTP/1.1 $message");
  } else {
    // それ以外の場合は500 Internal Server Errorとしてメッセージを出力する
    header('500 Internal Server Error');
    header('Content-Type: text/plain');
    echo $message;
  }
}