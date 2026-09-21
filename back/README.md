# メシイコ バックエンド

Go + Gin + GORM + SQLite のAPIサーバ．クリーンアーキテクチャの構成で，
OpenAPIのスキーマからコードを生成する．


## 必要なもの

- Docker / Docker Compose（Docker Desktop でよい）

## 起動

```bash
make up
```

| URL | 中身 |
| --- | --- |
| http://localhost:8080/api/v1 | API |
| http://localhost:8081 | Swagger UI（「Try it out」から直接APIを叩ける） |

動作確認：

```bash
curl http://localhost:8080/api/v1/health
curl -X POST http://localhost:8080/api/v1/users -H 'Content-Type: application/json' -d '{"name":"nishi"}'
curl http://localhost:8080/api/v1/users/1 -H 'Authorization: Bearer any-token'
```

## いまはモックAPI

**`/health` 以外のエンドポイントはまだロジックが入っておらず，固定のデータを返す．**
DBの読み書きもしない．フロント・ウィジェットの開発を先に進められるようにするためで，
レスポンスの形は `api/openapi.yaml` のとおり．

認証も本物ではない．

| 種類 | ヘッダ | モックでの扱い |
| --- | --- | --- |
| 利用者 | `Authorization: Bearer <任意の文字列>` | 中身は検証しない．空だと401 |
| 管理者（Cron用） | `X-Admin-Token: dev-admin-token` | `ADMIN_TOKEN` と一致するかだけ見る |

モックが返すもの：

- `GET /answers` は「行ける」「行けない」「未回答」が1人ずつ入った3件を返す
- `PUT /answers/me` は送った内容をそのまま返す（保存はしないので `GET /answers` は変わらない）
- `GET /polls/{date}` は常に `open`
- `GET /users/{userId}` は id が 1〜3 のときだけ 200．それ以外は 404

```bash
# 回答一覧（3状態がすべて入っている）
curl 'http://localhost:8080/api/v1/answers?date=2026-09-20' -H 'Authorization: Bearer any-token'

# 自分の回答を更新
curl -X PUT http://localhost:8080/api/v1/answers/me \
  -H 'Authorization: Bearer any-token' -H 'Content-Type: application/json' \
  -d '{"status":"available","timeSlots":["18:00","18:30"]}'

# 管理用（Cronから叩く想定）
curl -X PUT http://localhost:8080/api/v1/polls/2026-09-20 \
  -H 'X-Admin-Token: dev-admin-token' -H 'Content-Type: application/json' \
  -d '{"status":"closed"}'
```

固定データは `internal/interfaces/handler/mock.go` にまとめてある．
ロジックを実装するときは，各ハンドラの `mock〜` 呼び出しを usecase の呼び出しに
置き換えていけばよく，最後に `mock.go` を消せば終わる．

ソースを保存すると [air](https://github.com/air-verse/air) が自動でビルドし直す．

## よく使うコマンド

```bash
make help       # 一覧を表示
make up         # 起動（API + Swagger UI）
make down       # 停止
make clean      # 停止してDBのボリュームも削除
make logs       # APIのログを追う
make test       # テスト
make lint       # gofmt + go vet
make generate   # openapi.yaml から Go のコードを再生成
make sh         # 開発コンテナのシェルに入る
make prod       # 本番相当（ビルド済みバイナリ）で起動
```

## ディレクトリ構成

```
back/
├── api/
│   ├── openapi.yaml          ← APIの正（Single Source of Truth）
│   └── oapi-codegen.yaml     コード生成の設定
├── cmd/api/main.go           起動処理だけ
└── internal/
    ├── config/               環境変数の読み込み
    ├── domain/               ← 何にも依存しない中心
    │   ├── model/              エンティティ（User など）
    │   ├── repository/         永続化のインタフェース（実装は持たない）
    │   └── apperr/             レイヤー共通のエラー
    ├── usecase/              ← アプリケーションの手順．Gin も GORM も知らない
    ├── interfaces/           ← 外の世界との変換
    │   ├── handler/            HTTP ↔ usecase の変換
    │   ├── middleware/         ログ・CORS・認証
    │   └── router/             Ginの組み立て
    ├── infrastructure/       ← 技術的な詳細
    │   ├── database/           SQLite接続・マイグレーション
    │   └── persistence/        repository のGORM実装
    ├── registry/             依存の組み立て（DI）を集約
    └── generated/openapi/    ← 自動生成．手で編集しない
```

### 依存の向き

```
interfaces ──┐
             ├──> usecase ──> domain（model / repository / apperr）
infrastructure ┘                 ↑
                                 └── infrastructure が repository を実装する
```

内側（domain）は外側を一切importしない．この向きを守っている限り，
DBをSQLiteからCloudflare D1に変えても，変更は `infrastructure/persistence` で済む．

## 機能の追加手順

例として「お店」機能を足す場合：

1. **`api/openapi.yaml` にエンドポイントとスキーマを書く**（ここが出発点）
2. `make generate` で `internal/generated/openapi` を再生成する
3. `internal/domain/model/shop.go` にエンティティを書く
4. `internal/domain/repository/shop.go` にインタフェースを書く
5. `internal/infrastructure/persistence/` に GORM 実装と `*Record` を足し，
   `model.go` の `AllModels()` にテーブルを登録する
6. `internal/usecase/shop.go` に手順を書く
7. `internal/interfaces/handler/shop.go` にハンドラを書き，
   `handler.go` の `Handler` 構造体に埋め込む
8. `internal/registry/registry.go` に3行（repo / usecase / handler）足す
9. `make test`

`Handler` が `openapi.StrictServerInterface` を満たしていないと**コンパイルが通らない**ので，
openapi.yaml に書いたエンドポイントの実装漏れは必ずビルド時に見つかる．

## 設定（環境変数）

| 変数 | 既定値 | 説明 |
| --- | --- | --- |
| `PORT` | `8080` | 待ち受けポート |
| `GIN_MODE` | `debug` | `debug` / `release` / `test` |
| `DB_PATH` | `./data/meshi-iko.db` | SQLiteのファイル．`:memory:` でインメモリ |
| `CORS_ALLOWED_ORIGIN` | `*` | 許可するOrigin |
| `SHUTDOWN_TIMEOUT_SECONDS` | `10` | 終了時に処理中のリクエストを待つ秒数 |

compose.yaml で既定値を渡しているので，通常は何も設定しなくてよい．
上書きしたいときは `cp .env.example .env` する．

## 設計上の決めごと

- **スキーマ駆動**：`api/openapi.yaml` が正．Goの型・Ginのルーティングはそこから生成する．
  手で書いたハンドラと仕様がズレたらコンパイルエラーになる．
- **SQLiteドライバは純Go実装**（`glebarez/sqlite`）．CGOが要らないので，
  Dockerイメージが小さく，ビルドがホスト環境に左右されない．
- **AutoMigrate を使っている**．カラム削除やリネームは反映されないので，
  そこまで必要になったらマイグレーションツールを入れること．
- **エラーは `apperr.Kind` で表現する**．ユースケース層はHTTPステータスを知らず，
  handler が `Kind` をステータスに翻訳する．
- **dev と prod でComposeのプロジェクト名を分けている**．
  devはroot，prodは非rootで動くため，DBのボリュームを共有すると権限が衝突する．

## 未対応（デザインドック参照）

この枠組みには health と users（サンプル）しか入っていない．
`polls` / `answers` / FCMトークン登録 / Cronからの配信・締め切り / Slack連携は
`docs/DesignDoc.md` を見て上の手順で足していく．

端末トークン認証のミドルウェアは `internal/interfaces/middleware/auth.go` に
用意してあるが，まだどのルートにも適用していない．
