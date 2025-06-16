# PHPWebP Concept - PHP による WebP 対応の 自動化

**Web サイトの画像を次世代フォーマット WebP で軽量化し、適切に配信する運用を PHP で自動化** するコンセプト実装です。

最もポピュラーな Web サーバー構成である Apache + PHP 向けのプログラムで、LP や個人の小規模サイトでの利用を想定しています。

## 機能

従来フォーマット ⇄ WebP の双方向変換と、ブラウザの対応状況に応じた配信機能を提供します。

- WebP 対応ブラウザで従来フォーマット(JPEG/PNG/GIF)の画像をリクエスト → 軽量な WebP に変換して配信。
- WebP 非対応ブラウザで WebP 画像をリクエストすると PNG に変換して閲覧可能に。

## 動作環境

- Linux x86_64
  - 他の環境も `libwebp` の追加インストールにより対応可能
- Apache 2.x
  - 要 Rewrite モジュール
  - .htaccess または conf ファイルによる設定変更ができること
- PHP
  - 対応バージョンは `.github/workflows/cicd.yml` で管理
  - 現在対応: 5.6 / 7.0 / 7.1 / 7.2 / 7.3 / 7.4 / 8.0 / 8.1 / 8.2

## インストール

1. Releases から最新版の ZIP ファイルをダウンロードして展開してください。
2. `phpwebp-concept`ディレクトリを Web サーバーのドキュメントルートにアップロードします。
3. `.htaccess-example`を参考に、ドキュメントルートまたは画像の最上位ディレクトリに`.htaccess`ファイルを作成します。

すでに `.htaccess` が存在する場合は、既存の設定を損ねないように記述をマージしてください。

### OS が Linux x86_64 以外の場合

次の PHP プログラムを実行し、ディレクトリ名を確認します(例: `winnt-amd64`)。

```php
<?php echo strtolower(PHP_OS . '-' . php_uname('m'));
```

CLI であれば次のコマンドでも代替可能です。

```bash
php -r "echo strtolower(PHP_OS . '-' . php_uname('m'));"
```

`phpwebp-concept/bin` 以下に、上記名称のディレクトリを作成し、[libwebp](https://developers.google.com/speed/webp/download) に同梱される次のプログラムファイルを`実行可能なファイル`として配置してください。

- cwebp
- dwebp
- webpinfo
- gif2webp

## 動作確認

ブラウザの開発者ツールの `Network`タブ等で確認してください。

### 従来フォーマットから WebP への変換

WebP 対応ブラウザの開発者ツール `Network` タブ等で確認してください。

1. 従来フォーマット(JPEG/PNG/GIF)の画像レスポンスの`Content-Type`ヘッダが`image/webp`になっていること。
2. 従来フォーマット(JPEG/PNG/GIF)の画像レスポンスのデータサイズや`Content-Length`ヘッダが元画像より小さくなっていること。

### WebP から PNG への変換

WebP 非対応ブラウザで開発者ツール `Network` タブなどとあわせて確認してください。

1. WebP 画像を閲覧できること。
2. WebP 画像レスポンスの`Content-Type`ヘッダが`image/png`になっていること。

## ライセンス

このプログラムは MIT ライセンスに同意して利用ください。

## カスタマイズ / コントリビュート

Docker および Go 1.21以上 を事前にインストールし、本プロジェクトを `clone` してください。

### テスト

Go言語で実装されたテストスイートを使用します。dockertestを使用してPHPコンテナを動的に起動し、テストを実行します。

```bash
# 単一のPHPバージョンをテスト（デフォルト: 8.1）
make test

# 特定のPHPバージョンをテスト
PHP_VERSION=7.4 make test

# すべてのPHPバージョンを順次テスト
make test-all
```

### ビルド

次のコマンドでリリースパッケージをビルドします。

```bash
# デフォルトバージョン（v1.0.0）でビルド
make build

# 特定のバージョンでビルド
VERSION=v1.0.1 make build
```

リリースパッケージは `built/` ディレクトリに生成されます。

### クリーンアップ

```bash
# テストキャッシュとDockerコンテナをクリーンアップ
make clean
```
