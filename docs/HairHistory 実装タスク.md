# HairHistory 実装タスク


---

## MVP（完了）

- [x] 匿名ユーザーID発行（`POST /api/users`）
- [x] 履歴一覧（`GET /api/users/{userId}/histories`）
- [x] 履歴作成（`POST /api/users/{userId}/histories`）
- [x] 履歴更新（`PUT /api/histories/{historyId}`）
- [x] 履歴削除（`DELETE /api/histories/{historyId}`）
- [x] QR共有（URLをQR表示）
- [x] ランディングUI
- [x] アーカイブUI
- [x] APIエラー表示（フロント）

## 優先度高（次にやる）

- [x] 作成系 API を `201 Created` にする
- [x] 本番向けに 5xx エラー文言を固定化する
- [x] `main.go` の起動方式を一本化（`signal.NotifyContext`）

## 優先度中（その次）

- [x] パーマ要件と実装の整合（`ServiceType` の扱いを決定）
- [x] 施術の複数選択UI
- [x] ルート厳密化（将来衝突回避）
- [x] `usecase/request` の `*http.Request` 依存整理
- [x] PRIVACY / TERMS 実リンク化（お問い合わせページは設けない）
- [x] ルータ初期化は `chi.NewRouter()` に統一（`uchi` ラッパー削除）
- [x] CORS デフォルト時の警告ログ（`HAIR_CORS_ORIGINS` 未設定、またはパース後に有効オリジンなし）— `controller.ResolveCORSOrigins`

## レビューメモ（記録のみ・依頼時点ではコード対応なし）

### LGTM

- `apps/main/app/controller/*.go` — レイヤー責務の分離として良い変更
- `apps/main/app/usecase/request/*.go` — 純粋な DTO + バリデーションになりテスタビリティ向上
- `uchi.go` 削除 — 薄いラッパーの削除は妥当
- `apps/web/app/privacy/page.tsx`, `apps/web/app/terms/page.tsx` — Metadata 付き静的ページとして適切
- `docs/` — コード変更と CI の `npm install --no-audit --no-fund` / `npm run build` の記述が整合

### nits

- `.github/workflows/ci.yml` — ファイル末尾改行なしの指摘あり（@ko-tarou / POSIX）。**対応済み**（末尾 newline を追加）
- `cache: npm` は `package-lock.json` の存在を前提とする。現状リポジトリ直下に `package-lock.json` あり（確認済み）
- `hair_history_create.go`（controller）— `Decode` 失敗時はエラー文字列がそのままクライアントに返る。種別区別はないが現段階では 400 で十分。将来メッセージをラップするか検討してよい。
- `apps/web/app/page.tsx` — CONTACT リンク削除はタスクドキュメントと整合しており問題なし
- `privacy/page.tsx`（および同様の法務ページ）— 「最終更新」日付がハードコード。運用で更新忘れに注意。現段階では許容範囲。

### Q（未決・記録）

- `users_create.go`（controller）— `request.NewCreateUser()` が常に `nil` error を返す。将来バリデーション追加予定か。不要なら `error` を外してシンプルにする選択肢もある（@ko-tarou）。**答え（記録）**: 将来拡張を見越して `error` を残す意図であれば一貫してよい。

### その他（過去レビューからのメモ）

- `CreateHistory` で URL から `UserID` をセット後に `Decode(in)` するが、body の `UserID` 相当キーで上書きされないか — **`UserID` は `json:"-"` のため JSON からは設定されない**（`UpdateHistory.HistoryID` も同様）。別名キーで別フィールドにマップされるタグがない限り問題なし。

---

## 拡張（任意）

- 認証（ログイン）
- 分析・予約・設定
- ダメージサマリー
- 通知 / リマインド
- E2Eテスト整備
- docs の作業メモ運用ルール決定（`docs/` 管理 or `.gitignore`）
- CSSエラー色を変数化（例: `--color-error`）
- `apps/web/app/h/[userId]/page.tsx` の state 整理（将来的に `useReducer` やカスタムフックへ分離）

---

更新ルール:

- 完了した項目はチェックを付ける
- 仕様変更は `要件定義` / `基本設計` へ反映する
