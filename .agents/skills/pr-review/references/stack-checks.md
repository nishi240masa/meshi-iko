# スタック別の観点

該当するファイルが変更された節だけ読めばよい。見出しがそのまま観点 ID になる。

- `go-gin-gorm`: `back/**/*.go`
- `db-neon-lambda`: DB接続・モデル・クエリ・Dockerfile・Lambda 関連
- `scheduler-fcm`: 管理用API、通知、devices
- `react-ts`: `front/src/**`
- `tauri-pwa`: `front/src-tauri/**`、PWA の manifest / Service Worker
- `widget`: Glance / Scriptable
- `ci-tooling`: `.github/**`、`.husky/**`、Dockerfile、compose

## go-gin-gorm
- `context.Context` をハンドラからDBまで伝えているか（GORM は `WithContext`）。
- goroutine のリーク・データ競合。レスポンス後のバックグラウンド処理は Lambda では凍結されるので前提にできない。
- GORM: `First` と `Find` の使い分け（`ErrRecordNotFound`）、`Preload`・`Joins` での N+1 回避、トランザクションの範囲、構造体での `Updates` はゼロ値を無視する、本番起動時に `AutoMigrate` を走らせていないか（`cmd/migrate` から流す設計）。
- oapi-codegen の生成型・インタフェースを使っているか、手書きで型を重複定義していないか、spec の制約（required、enum、format、上限）と実装の検証が一致しているか。

## db-neon-lambda
- SQLite 固有の振る舞いに依存していないか（型の緩さ、LIKE の大文字小文字、日時の保存形式、ロック、方言の違うSQL）。テストは Postgres で実行する設計。
- Neon の pooler（PgBouncer transaction mode）で使えない機能に依存していないか: `SET`（`search_path` を含む）、`LISTEN/NOTIFY`、SQL レベルの `PREPARE`、セッション単位の advisory lock、トランザクションをまたぐ一時テーブル。マイグレーションは直接接続で行う設計か。
- Lambda: DB接続・クライアントはハンドラの外で1回だけ初期化しているか、`MaxOpenConns` を絞っているか、インメモリのキャッシュやレート制限をインスタンス間で共有できる前提にしていないか。

## scheduler-fcm
- API Destination から呼ばれるエンドポイントが `X-Admin-Token` で認証され、再送で二重実行されても安全（冪等）か。
- API Destination は5秒でタイムアウトし、再送は2回・3600秒に絞る設計（DesignDoc §6）。処理がこれに収まり、再送で翌日分を締め切らないか。
- FCM は登録トークンではなく FID を使う。フロントで `getToken()` を呼んでいないか、`onRegistered()` の登録が1か所か。
- 送信は `SendEachForMulticast`（500件単位）、`UNREGISTERED`/`INVALID_ARGUMENT` で devices の行を消しているか。ペイロードに個人情報や秘密を載せていないか。

## react-ts
- openapi-fetch と生成型（`front/src/generated/api.ts`）を使っているか（`any` や型アサーションで逃げていないか）。エラー分岐を処理しているか。
- Zod で信頼境界（フォーム入力、URL、ローカルストレージ）を検証しているか。同じスキーマを二重に定義していないか。
- `useEffect` の乱用（派生値を state に持つ、イベントハンドラで済む処理を Effect で行う）、依存配列、古いリクエストの結果で新しい結果を上書きする競合。
- React Compiler は未使用なので、重い計算や長いリストでの `useMemo` / `useCallback` の要否を判断する（設定が変わっていないか `front/vite.config.ts` で確認）。
- TanStack Router: パス・検索パラメータの型と検証、認証が必要なルートの保護。
- トークンの保存場所が設計ドキュメントどおりか（PWA と Tauri で異なりうる）。

## tauri-pwa
- `front/src-tauri/capabilities/default.json` が最小権限か（ワイルドカード、広すぎるファイルシステムのスコープ、不要な shell 権限）。`front/src-tauri/tauri.conf.json` に CSP が設定されているか。
- Service Worker が認証付きAPIレスポンスを不適切にキャッシュしていないか、更新時に古いキャッシュを破棄できるか。
- iPhone（PWA）と Android（Tauri）で使えるAPIが違う箇所に両方の分岐があるか。

## widget
- 更新頻度とバッテリー・APIのレート、認証情報の保持方法、APIスキーマ変更への追従。
- Android は Kotlin Glance、iOS は Scriptable。配布する JS スクリプトにトークンや秘密を埋め込んでいないか。

## ci-tooling
- GitHub Actions（backend.yml / frontend.yml の Rust ジョブを含む）の `permissions` が最小か、サードパーティ Action をコミットSHAで固定しているか、`pull_request_target` を危険な形で使っていないか、secrets が露出しないか。
- 生成物の一致チェックなど、既存のCIの確認を弱めていないか。`.husky` の pre-commit を弱めていないか。
- Dockerfile / compose: 不要なファイルの同梱、root での実行、Lambda Web Adapter の設定。
