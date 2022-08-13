<?php

require(__DIR__ . '/common.php');

/**
 * ApacheのRewriteモジュールによってルーティングされた画像へのリクエストを処理する
 * 
 * - リクエストパスから実画像ファイルを特定し、WebPキャッシュがなければ変換する
 * - キャッシュをWebPファイルとして返すが、キャッシュファイルのサイズが0バイトの場合は変換失敗として元画像を返す
 */
function main() {
  // リクエストを解析する
  // $image_file_pathには従来フォーマット、$cache_file_pathには対応するWebPキャッシュのパスが入る
  parse_request($method, $image_file_path, $cache_file_path, $cache_key);

  // キャッシュファイルがまだなければ画像をWebPに変換してキャッシュする
  if (!file_exists($cache_file_path)) {
    try {
      convert_to_webp($image_file_path, $cache_file_path);
    } catch(\Exception $ex) {
      // 変換エラーの原因はログに出力する
      error_log("Failed to convert to WebP: $image_file_path; Message: " . $ex->getMessage());
      // 変換に失敗した場合はキャッシュファイルを0バイトにする
      // これは変換を繰り返さず、かつ変換の失敗を識別するため
      file_put_contents($cache_file_path, '');
    }
  }

  // キャッシュファイルの実体があればWebPとして返す
  if (file_exists($cache_file_path) && ($cache_size = filesize($cache_file_path)) > 0) {
    header('Content-Type: image/webp');
    header('Content-Length: ' . filesize($cache_file_path));
    header('X-Cache-Key: ' . $cache_key);

    // 軽量化の結果をヘッダに出力する
    $image_size = filesize($image_file_path);
    header('X-PHPWebP-Stats: ' . sprintf('status=success; original=%0.1fkb; ratio=%0.2f%%;', $image_size / 1024, $cache_size * 100 / $image_size));

    if ($method !== 'HEAD') readfile($cache_file_path);
  } else {
    // キャッシュファイルがない、またはファイルサイズが0なら元画像を返す
    header('Content-Type: ' . mime_content_type($image_file_path));
    header('Content-Length: ' . filesize($image_file_path));

    // 軽量化の失敗をヘッダに出力する
    header('X-PHPWebP-Stats: status=failure;');

    if ($method !== 'HEAD') readfile($image_file_path);
  }
}

/**
 * ファイルのフォーマットに応じてWebPに変換する
 */
function convert_to_webp($input_image_path, $output_webp_path) {
  // 画像のフォーマットを識別
  $mime = mime_content_type($input_image_path);

  if ($mime === 'image/jpeg') {
    safe_webp_exec_via_stdin('cwebp', "-q 80 -o $output_webp_path -- -", $input_image_path);
  } else if ($mime === 'image/png') {
    safe_webp_exec_via_stdin('cwebp', "-lossless -o $output_webp_path -- -", $input_image_path);
  } else if ($mime === 'image/gif') {
    safe_webp_exec_via_stdin('gif2webp', "-o $output_webp_path -- -", $input_image_path);
  } else {
    throw new \Exception('Unsupported image type');
  }

  // 前後のファイルサイズをチェックする
  // WebPが元のファイルサイズを上回ることが稀にある(GIFは比較的多い)ため
  if (filesize($input_image_path) < filesize($output_webp_path)) {
    throw new \Exception('WebP got larger than original image');
  }
}

try {
  // メイン処理
  main();
} catch(\Exception $ex) {
  handle_as_http_exception($ex);
}
