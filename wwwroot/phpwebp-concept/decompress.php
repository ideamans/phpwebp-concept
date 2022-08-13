<?php

require(__DIR__ . '/common.php');

/**
 * ApacheのRewriteモジュールによってルーティングされた画像へのリクエストを処理し、WebPを従来フォーマットに変換する
 * 
 * - リクエストパスから実画像ファイルを特定し、従来フォーマットのキャッシュがなければ変換する
 * - キャッシュを従来フォーマットのファイルとして返すが、キャッシュファイルのサイズが0バイトの場合は変換失敗として元画像を返す
 */
function main() {
  // リクエストを解析する
  // $image_file_pathにはWebP、$cache_file_pathには対応する従来フォーマットのキャッシュが入る
  parse_request($method, $image_file_path, $cache_file_path, $cache_key);

  // キャッシュファイルがまだなければ画像をWebPに変換してキャッシュする
  if (!file_exists($cache_file_path)) {
    try {
      convert_from_webp($image_file_path, $cache_file_path);
    } catch(\Exception $ex) {
      // 変換エラーの原因はログに出力する
      error_log("Failed to convert from WebP: $image_file_path; Message: " . $ex->getMessage());
      // 変換に失敗した場合はキャッシュファイルを0バイトにする
      // これは変換を繰り返さず、かつ変換の失敗を識別するため
      file_put_contents($cache_file_path, '');
    }
  }

  // キャッシュファイルの実体があれば従来フォーマットとして返す
  if (file_exists($cache_file_path) && ($cache_size = filesize($cache_file_path)) > 0) {
    header('Content-Type: ' . mime_content_type($cache_file_path));
    header('Content-Length: ' . filesize($cache_file_path));
    header('X-Cache-Key: ' . $cache_key);

    // 変換によるサイズの増大結果をヘッダに出力する
    $image_size = filesize($image_file_path);
    header('X-PHPWebP-Stats: ' . sprintf('status=success; webp=%0.1fkb; ratio=%0.2f%%;', $image_size / 1024, $cache_size * 100 / $image_size));

    if ($method !== 'HEAD') readfile($cache_file_path);
  } else {
    // キャッシュファイルがない、またはファイルサイズが0なら元の画像(推定WebP)を返す
    $mimetype = mime_content_type($image_file_path);
    header('Content-Type: ' . ($mimetype ? $mimetype : 'image/webp'));
    header('Content-Length: ' . filesize($image_file_path));

    // 変換の失敗をヘッダに出力する
    header('X-PHPWebP-Stats: status=failure;');

    if ($method !== 'HEAD') readfile($image_file_path);
  }
}

/**
 * ファイルのフォーマットに応じてWebPに変換する
 */
function convert_from_webp($input_webp_path, $output_image_path) {
  // dwebpでPNGに変換する
  // LossyでAlphaがない場合はJPEGにしたいところだが、
  // dwebpはJPEGに変換する機能がないためPNGのみサポートする
  safe_webp_exec_via_stdin('dwebp', "-o $output_image_path -- -", $input_webp_path);
}

try {
  // メイン処理
  main();
} catch(\Exception $ex) {
  handle_as_http_exception($ex);
}
