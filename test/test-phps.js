import test from 'ava'
import Axios from 'axios'
import ImageType from 'image-type'

// テストするPHPバージョン
const phps = [
  { version: '5.4', dockerHost: 'php-54', localHost: 'localhost:10540' },
  { version: '5.5', dockerHost: 'php-55', localHost: 'localhost:10550' },
  { version: '5.6', dockerHost: 'php-56', localHost: 'localhost:10560' },
  { version: '7.0', dockerHost: 'php-70', localHost: 'localhost:10700' },
  { version: '7.1', dockerHost: 'php-71', localHost: 'localhost:10710' },
  { version: '7.2', dockerHost: 'php-72', localHost: 'localhost:10720' },
  { version: '7.3', dockerHost: 'php-73', localHost: 'localhost:10730' },
  { version: '7.4', dockerHost: 'php-74', localHost: 'localhost:10740' },
  { version: '8.0', dockerHost: 'php-80', localHost: 'localhost:10800' },
  { version: '8.1', dockerHost: 'php-81', localHost: 'localhost:10810' },
]

const targetPhp = process.env.PHP
const targetPhps = phps.filter((php) => !targetPhp || php.version === targetPhp)

for (const php of targetPhps) {
  test(`動作確認: PHP ${php.version}`, async (t) => {
    // WebP対応エージェント(Accept: image/webp あり)のテスト
    // 従来フォーマットからWebPに正常に変換されるケース
    await Promise.all(
      [
        '/testing/regular.jpg',
        '/testing/regular.png',
        '/testing/animation.gif',
      ].map((path) =>
        testValidImageRequest(
          t,
          'WebP対応エージェント',
          'image/webp,*/*',
          php,
          path,
          200,
          'image/webp'
        )
      )
    )

    // WebPに正常に変換されないケース(非対応のBMP)
    await Promise.all(
      ['/testing/bmp.jpg'].map((path) => {
        return testInvalidRequest(
          t,
          'WebP対応エージェント',
          'image/webp,*/*',
          php,
          path,
          200,
          'image/x-ms-bmp'
        )
      })
    )

    // WebPに正常に変換されないケース(非対応のCMYK画像)
    await Promise.all(
      ['/testing/cmyk.jpg'].map((path) =>
        testValidImageRequest(
          t,
          'WebP対応エージェント',
          'image/webp,*/*',
          php,
          path,
          200,
          'image/jpeg'
        )
      )
    )

    // WebPに直接アクセスするケース
    await Promise.all(
      ['/testing/lossy.webp', '/testing/lossless.webp'].map((path) =>
        testValidImageRequest(
          t,
          'WebP対応エージェント',
          'image/webp,*/*',
          php,
          path,
          200,
          'image/webp'
        )
      )
    )

    // WebP非対応エージェント(Accept: image/webp なし)のテスト
    // 従来フォーマットの画像にアクセスするケース
    await Promise.all(
      [
        ['/testing/regular.jpg', 'image/jpeg'],
        ['/testing/regular.png', 'image/png'],
        ['/testing/animation.gif', 'image/gif'],
      ].map(([path, mimeType]) =>
        testValidImageRequest(
          t,
          'WebP非対応エージェント',
          '*/*',
          php,
          path,
          200,
          mimeType
        )
      )
    )

    // WebP画像にアクセスするケース
    await Promise.all(
      [
        ['/testing/lossy.webp', 'image/png'],
        ['/testing/lossless.webp', 'image/png'],
      ].map(([path, mimeType]) =>
        testValidImageRequest(
          t,
          'WebP非対応エージェント',
          '*/*',
          php,
          path,
          200,
          mimeType
        )
      )
    )

    // WebP非対応エージェントとしてアクセスするがPNGに変換できないケース(実体がBMP)
    const bmpMimeType = ['5.4', '5.5'].includes(php.version)
      ? 'application/octet-stream'
      : 'image/x-ms-bmp'
    await Promise.all(
      [['/testing/bmp.webp', bmpMimeType, 'image/bmp']].map(
        ([path, mimeType, dataMimeType]) =>
          testValidImageRequest(
            t,
            'WebP非対応エージェント',
            '*/*',
            php,
            path,
            200,
            mimeType,
            dataMimeType
          )
      )
    )

    // 404 Not Foundのテスト
    await Promise.all(
      [
        ['WebP対応エージェント', 'image/webp,*/*'],
        ['WebP非対応エージェント', '*/*'],
      ].map(([agentLabel, acceptHeader]) =>
        testInvalidRequest(
          t,
          agentLabel,
          acceptHeader,
          php,
          '/testing/notfound.jpg',
          404
        )
      )
    )
  })
}

/**
 * ローカルホストWebサーバへの有効な画像リクエストをテスト
 *
 * @param {ava} t AVAのテストコンテキスト
 * @param {string} agentLabel エージェント名
 * @param {string} acceptHeader Acceptリクエストヘッダ値
 * @param {php} php PHPバージョン情報
 * @param {string} path URLパス
 * @param {number} expectedStatus 期待するステータスコード
 * @param {string} expectedMimeType 期待するMIMEタイプ
 */
async function testValidImageRequest(
  t,
  agentLabel,
  acceptHeader,
  php,
  path,
  expectedStatus,
  expectedMimeType,
  expectedDataMimeType
) {
  expectedDataMimeType ||= expectedMimeType
  const host =
    process.env.PHP_HOST === 'docker' ? php.dockerHost : php.localHost
  const res = await Axios.get(`http://${host}${path}`, {
    headers: { Accept: acceptHeader },
    responseType: 'arraybuffer',
    validateStatus: () => true,
  })

  t.is(
    res.status,
    expectedStatus,
    `${agentLabel}: ${path} が ステータスコード ${expectedStatus} を返すこと`
  )

  t.is(
    res.headers['content-type'],
    expectedMimeType,
    `${agentLabel}: ${path} がヘッダ "Content-Type: ${expectedMimeType}" を返すこと`
  )

  const imageType = await ImageType(res.data)
  t.is(
    imageType.mime,
    expectedDataMimeType,
    `${agentLabel}: ${path} の実体が ${expectedDataMimeType} であること`
  )
}

/**
 * ローカルホストWebサーバーへの無効なリクエストをテスト
 *
 * @param {ava} t AVAのテストコンテキスト
 * @param {string} agentLabel エージェント名
 * @param {string} acceptHeader Acceptリクエストヘッダ値
 * @param {php} php PHPバージョン情報
 * @param {string} path URLパス
 * @param {number} expectedStatus 期待するステータスコード
 */
async function testInvalidRequest(
  t,
  agentLabel,
  acceptHeader,
  php,
  path,
  expectedStatus
) {
  const host =
    process.env.PHP_HOST === 'docker' ? php.dockerHost : php.localHost
  const res = await Axios.get(`http://${host}${path}`, {
    headers: { Accept: acceptHeader },
    responseType: 'arraybuffer',
    validateStatus: () => true,
  })

  t.is(
    res.status,
    expectedStatus,
    `${agentLabel}: ${path} が ステータスコード ${expectedStatus} を返すこと`
  )
}
