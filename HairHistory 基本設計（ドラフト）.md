# HairHistory 基本設計（ドラフト）

要件定義（`要件定義.md`）をもとにした、画面・データ・API のざっくり設計です。

目的は **「髪の履歴をメモして、ダメージを美容師が把握しやすくする」** こと。

- --

## 1. 画面構成（MVP）

1. **トップ / 初回アクセス画面**

- 説明テキスト（このアプリの目的）

- 「はじめる」ボタン → ユーザーID（URL用ID）を自動発行し、自分用の履歴ページへ遷移

2. **マイ履歴一覧画面（/h/:userId）**

- ユーザーURLでアクセスするメイン画面

- 表示内容（例）:

- ヘッダー: プロジェクト名、簡単な説明（「美容師さんと一緒にダメージを確認できます」など）

- 「新しい履歴を追加」ボタン

- 履歴一覧（新しい順）

- 日付

- 施術内容（例: カラー / ブリーチ / 縮毛矯正）

- サロン名・スタイリスト名（簡易表示）

- ダメージが分かるような短いサマリー（例: 「ブリーチ2回目 / 毛先ダメージ強め」など）

3. **履歴入力・編集画面（モーダル or 別ページ）**

- 1件の履歴を登録・編集する画面

- 入力項目（要件から）:

- 日付

- 施術内容

- チェックボックス例: カラー / ブリーチ / 縮毛矯正 / トリートメント / その他（自由入力）

- サロン名（フリーテキスト）

- スタイリスト名（フリーテキスト）

- メモ

- ダメージに関するメモを書きやすいように、プレースホルダー例:

- 「例: 毛先がかなりダメージ / ブリーチ〇回目 / 薬剤がしみた など」

4. **共有ビュー（美容師さんが見る想定の画面）**

- 基本はマイ履歴一覧画面と同一URL（/h/:userId）

- スマホを見せる / URL を送る前提なので、

- 直近3〜5件の履歴を上にまとめて表示

- 各履歴に「施術内容」と「ダメージ関連メモ」が一目で分かるように

- --

## 2. データ設計（詳細）

### 2.1 テーブル: users（匿名ユーザー）

| カラム名 | 型 | NULL | 説明 |
|---------|----|------|------|

| id | TEXT | ❌ | ユーザーID（URL用ID、nanoid等） |

| name | TEXT | ✅ | 表示名（「〇〇さんの履歴」用） |

| email | TEXT | ✅ | メールアドレス（将来ログイン用） |

| password_hash | TEXT | ✅ | パスワードハッシュ（将来ログイン用） |

| last_login_at | TIMESTAMPTZ | ✅ | 最終ログイン日時 |

| is_deactivated | BOOLEAN | ❌ | 無効化フラグ（退会・BAN時に true、DEFAULT false） |

| created_at | TIMESTAMPTZ | ❌ | 作成日時 |

※ MVP では id / created_at のみ使用。name は任意入力。email 以降はログイン導入時に利用。

### 2.2 テーブル: hair_histories（髪の履歴）

PostgreSQL 想定の型で記載。

| カラム名 | 型 | NOT NULL | 説明 |
|---------|----|----------|------|

| id | UUID | ✅ | 履歴ID（サーバー側で生成） |

| user_id | UUID | ✅ | users.id への外部キー |

| date | DATE | ✅ | 施術日 |

| services | TEXT | ✅ | 施術内容。JSON 文字列（例: `["color","bleach"]`）として保存 |

| salon_name | TEXT | ❌ | サロン名（フリーテキスト） |

| stylist_name | TEXT | ❌ | スタイリスト名（フリーテキスト） |

| memo | TEXT | ❌ | メモ（ダメージ状況・薬剤・注意点など） |

| created_at | TIMESTAMPTZ | ✅ | 作成日時（サーバー側で設定） |

| updated_at | TIMESTAMPTZ | ✅ | 更新日時（サーバー側で設定） |

※ `services` は Go 側では `[]ServiceType` として扱い、DB には JSON 文字列で保存する想定。

- --

## 2.3 Go の型（domain / DTO）

### 2.3.1 施術内容の種類

```go

type ServiceType string

const (

ServiceTypeColor        ServiceType = "color"

ServiceTypeBleach       ServiceType = "bleach"

ServiceTypeStraightPerm ServiceType = "straight_perm" *// 縮毛矯正*

ServiceTypeTreatment    ServiceType = "treatment"

ServiceTypeOther        ServiceType = "other"

)

```

### 2.3.2 ドメインモデル

```go

type User struct {

ID           string     `json:"id"`

Name         string     `json:"name,omitempty"`

Email        string     `json:"email,omitempty"`

LastLoginAt   *time.Time `json:"lastLoginAt,omitempty"`

IsDeactivated bool       `json:"isDeactivated"`

CreatedAt    time.Time  `json:"createdAt"`

}

type HairHistory struct {

ID          string        `json:"id"`

UserID      string        `json:"userId"`

Date        time.Time     `json:"date"`

Services    []ServiceType `json:"services"`

SalonName   string        `json:"salonName"`

StylistName string        `json:"stylistName"`

Memo        string        `json:"memo"`

CreatedAt   time.Time     `json:"createdAt"`

UpdatedAt   time.Time     `json:"updatedAt"`

}

```

### 2.3.3 リクエスト / レスポンス用 DTO

```go

*// POST /api/users レスポンス*

type CreateUserResponse struct {

UserID string `json:"userId"`

}

*// GET /api/users/{userId}/histories レスポンス*

type ListHistoriesResponse struct {

Histories []HairHistory `json:"histories"`

}

*// POST /api/users/{userId}/histories リクエスト*

type CreateHistoryRequest struct {

Date        time.Time     `json:"date"`

Services    []ServiceType `json:"services"`

SalonName   string        `json:"salonName"`

StylistName string        `json:"stylistName"`

Memo        string        `json:"memo"`

}

*// POST /api/users/{userId}/histories レスポンス*

type CreateHistoryResponse struct {

History HairHistory `json:"history"`

}

*// PUT /api/histories/{historyId} リクエスト*

type UpdateHistoryRequest struct {

Date        *time.Time     `json:"date,omitempty"`

Services    *[]ServiceType `json:"services,omitempty"`

SalonName   *string        `json:"salonName,omitempty"`

StylistName *string        `json:"stylistName,omitempty"`

Memo        *string        `json:"memo,omitempty"`

}

*// PUT /api/histories/{historyId} レスポンス*

type UpdateHistoryResponse struct {

History HairHistory `json:"history"`

}

```

- --

## 3. API 設計（MVP 想定）

ベースURL: `/api`

1. **ユーザー作成**

- `POST /api/users`

- 内容: 新しい匿名ユーザーIDを発行

- レスポンス: `CreateUserResponse`

2. **履歴一覧取得**

- `GET /api/users/{userId}/histories`

- 役割: 指定ユーザーの履歴一覧を取得（新しい順）

- レスポンス: `ListHistoriesResponse`

3. **履歴作成**

- `POST /api/users/{userId}/histories`

- ボディ: `CreateHistoryRequest`

- レスポンス: `CreateHistoryResponse`

4. **履歴更新**

- `PUT /api/histories/{historyId}`

- ボディ: `UpdateHistoryRequest`

- レスポンス: `UpdateHistoryResponse`

5. **履歴削除（MVPで入れるかは検討）**

- `DELETE /api/histories/{historyId}`

- レスポンスボディなし or `{ "ok": true }`

- --

## 4. 「ダメージが分かりやすい」ためのUIポイントメモ

- 履歴一覧で、**直近のブリーチ・縮毛矯正の回数**がざっくり分かるようにしたい

- 例: 一覧の1行に「カラー / ブリーチ（通算3回目） / 縮毛矯正（前回から6ヶ月）」など

- メモ欄のプレースホルダーやラベルで、「ダメージ」「しみた」「切れ毛」などを書きやすくする

- 直近3〜5件をカード風に大きめ表示 → 美容師さんがパッと見て判断しやすい

- --

## 5. Go API のディレクトリ構成（案）

MVP向けのシンプル構成。`internal` 配下にレイヤーを分ける。

```text

backend/

cmd/

api/

main.go          // エントリポイント（HTTPサーバー起動）

internal/

http/

handler/

user_handler.go        // ユーザー関連ハンドラ（POST /api/users）

history_handler.go     // 履歴関連ハンドラ（GET/POST/PUT/DELETE）

router.go                // ルーティング定義（chi / mux / stdlib など）

domain/

model.go                 // User, HairHistory, ServiceType など

repository/

history_repository.go    // hair_histories 向けDBアクセス

user_repository.go       // users 向けDBアクセス

db.go                    // DB接続（PostgreSQL）初期化

service/

history_service.go       // 履歴のビジネスロジック（将来必要になったら）

config/

config.go                // 環境変数読み込み（DB URLなど）

```

- **最初は無理に分けすぎない方針**:

- まずは `handler` + `repository` + `domain` + `db` くらいの分割でOK。

- ロジックが増えてきたら、`service` レイヤーを厚くしていく。

- Next.js からは `/api/*` をこの Go API に向けて呼び出す想定（Railway 上でホスト）。

- --

この設計はドラフトです。画面項目やAPI・ディレクトリ構成の詳細は、実装しながら一緒に調整していきましょう。

