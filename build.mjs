#!/usr/bin/env zx

const Package = require('./package.json')

// バージョンの特定と確認
const packageVersion = `v${Package.version}`

// バージョンはコマンドライン引数を優先する
// ユースケースとしてはGitHub Actionでタグ名から渡す(例: v1.0.0)
const version = process.argv[3] || packageVersion
if (version !== packageVersion) {
  // 一致しない場合はリリースしない
  console.log(
    `version: ${version} is not match with package.json version: ${Package.version}`
  )
  process.exit(1)
}

// ディレクトリやファイルのパス
const release = `phpwebp-concept-${version}`
const workingDir = `built/${release}`
const builtZip = `built/${release}.zip`

// クリーンナップ
await $`rm -rf "${workingDir}" "${builtZip}"`
await $`mkdir -p "${workingDir}"`

// ZIPファイルの作成
await $`cp -a wwwroot/phpwebp-concept "${workingDir}/phpwebp-concept"`
await $`cp -a wwwroot/.htaccess "${workingDir}/htaccess-example.txt"`
await $`cd built && zip -r "${release}.zip" "${release}"`
await $`rm -rf "${workingDir}"`
