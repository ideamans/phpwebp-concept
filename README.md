# 🖼️ PHPWebP Concept - PHP による WebP 対応の 自動化

**Web サイトの画像を次世代フォーマット WebP で軽量化し、適切に配信する運用を PHP で自動化** するコンセプト実装です。

最もポピュラーな Web サーバー構成である Apache + PHP 向けのプログラムで、共用レンタルサーバーから VPS まで幅広い環境で利用可能です。

[![CI/CD](https://github.com/ideamans/phpwebp-concept/actions/workflows/cicd.yml/badge.svg)](https://github.com/ideamans/phpwebp-concept/actions/workflows/cicd.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## ✨ 機能

従来フォーマット ⇄ WebP の双方向変換と、ブラウザの対応状況に応じた配信機能を提供します。

- 🚀 WebP 対応ブラウザで従来フォーマット(JPEG/PNG/GIF)の画像をリクエスト → 軽量な WebP に変換して配信
- 🔄 WebP 非対応ブラウザで WebP 画像をリクエストすると PNG に変換して閲覧可能に
- ⚡ 変換結果をキャッシュして高速化
- 🔒 パストラバーサル攻撃を防ぐセキュアな実装
- 🏠 **レンタルサーバーでも動作** - root権限不要で共用サーバーでも利用可能

## 📋 動作環境

- Linux x86_64 / Linux ARM64 (aarch64)
  - 主要なレンタルサーバーで動作確認済み
  - 他の環境も `libwebp` の追加インストールにより対応可能
- Apache 2.x
  - 要 Rewrite モジュール
  - .htaccess または conf ファイルによる設定変更ができること
- PHP
  - 対応バージョンは `.github/workflows/cicd.yml` で管理
  - 現在対応: 5.6 / 7.0 / 7.1 / 7.2 / 7.3 / 7.4 / 8.0 / 8.1 / 8.2

## 📦 インストール

### 🏠 レンタルサーバーでの導入（推奨）

多くの共用レンタルサーバーで**追加料金なし**でWebP配信が可能になります：

1. [Releases](https://github.com/ideamans/phpwebp-concept/releases) から最新版の ZIP ファイルをダウンロード
2. FTPまたはファイルマネージャーで `phpwebp-concept` ディレクトリをアップロード
3. `.htaccess-example` を参考に `.htaccess` を設定

**対応確認済みのレンタルサーバー例：**
- さくらのレンタルサーバ
- エックスサーバー
- ロリポップ！
- その他 Apache + PHP が動作する環境

> 💡 **メリット**: root権限不要、追加モジュール不要、既存サイトに簡単導入可能

### 🔧 その他のアーキテクチャへの対応

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

## 🔍 動作確認

ブラウザの開発者ツールの `Network`タブ等で確認してください。

### 従来フォーマットから WebP への変換

WebP 対応ブラウザの開発者ツール `Network` タブ等で確認してください。

1. 従来フォーマット(JPEG/PNG/GIF)の画像レスポンスの`Content-Type`ヘッダが`image/webp`になっていること。
2. 従来フォーマット(JPEG/PNG/GIF)の画像レスポンスのデータサイズや`Content-Length`ヘッダが元画像より小さくなっていること。

### WebP から PNG への変換

WebP 非対応ブラウザで開発者ツール `Network` タブなどとあわせて確認してください。

1. WebP 画像を閲覧できること。
2. WebP 画像レスポンスの`Content-Type`ヘッダが`image/png`になっていること。

## 📄 ライセンス

このプログラムは MIT ライセンスに同意して利用ください。

## 💰 導入のメリット

### レンタルサーバーでのメリット
- **追加費用なし** - 既存のホスティングプランのまま画像を最適化
- **転送量削減** - 画像サイズを25-35%削減し、転送量制限を有効活用
- **表示速度向上** - Google PageSpeed Insights のスコア改善
- **SEO効果** - Core Web Vitals の改善により検索順位向上の可能性

## 🛠️ カスタマイズ / コントリビュート

### 必要な環境

- 🐳 Docker
- 🐹 Go 1.21以上
- 📝 Make

本プロジェクトを `clone` してください：

```bash
git clone https://github.com/ideamans/phpwebp-concept.git
cd phpwebp-concept
```

### 🧪 テスト

Go言語で実装されたテストスイートを使用します。dockertestを使用してPHPコンテナを動的に起動し、テストを実行します。

```bash
# 単一のPHPバージョンをテスト（デフォルト: 8.1）
make test

# 特定のPHPバージョンをテスト
PHP_VERSION=7.4 make test

# すべてのPHPバージョンを順次テスト
make test-all
```

### 🔨 ビルド

次のコマンドでリリースパッケージをビルドします。

```bash
# デフォルトバージョン（v1.0.0）でビルド
make build

# 特定のバージョンでビルド
VERSION=v1.0.1 make build
```

リリースパッケージは `built/` ディレクトリに生成されます。

### 🧹 クリーンアップ

```bash
# テストキャッシュとDockerコンテナをクリーンアップ
make clean
```

## 🤝 コントリビューション

### PHPバージョンの追加

1. [Docker Hub](https://hub.docker.com/_/php) で利用可能なPHPバージョンを確認
2. `.github/workflows/cicd.yml` のmatrixに追加
3. ローカルでテスト: `PHP_VERSION=X.X make test`
4. テストが通ったらPRを作成

### libwebpの更新

1. 現在のバージョンを確認: `cat .libwebp-version`
2. [Google Developers](https://developers.google.com/speed/webp/docs/precompiled) から新バージョンをダウンロード
3. バイナリを更新して `.libwebp-version` を更新
4. `make test-all` で全PHPバージョンでテスト

## 📊 パフォーマンス

- 画像サイズを平均25-35%削減（JPEG比）
- キャッシュにより2回目以降のアクセスは高速化
- WebPが元画像より大きい場合は元画像を配信

## 🔧 トラブルシューティング

### 画像が変換されない場合

1. Apache の mod_rewrite が有効か確認
2. `.htaccess` が正しく配置されているか確認
3. `phpwebp-concept/bin/` 内のバイナリに実行権限があるか確認

### エラーログの確認

PHPのエラーログまたはApacheのエラーログを確認してください。

## 🔗 関連リンク

- [WebP公式サイト](https://developers.google.com/speed/webp)
- [libwebpダウンロード](https://developers.google.com/speed/webp/docs/precompiled)
- [イシューレポート](https://github.com/ideamans/phpwebp-concept/issues)
