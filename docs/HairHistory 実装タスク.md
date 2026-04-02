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

- パーマ要件と実装の整合（`ServiceType` の扱いを決定）
- 施術の複数選択UI
- ルート厳密化（将来衝突回避）
- `usecase/request` の `*http.Request` 依存整理
- CONTACT / PRIVACY / TERMS 実リンク化
- `main.go` の CORS デフォルト設定に警告ログ追加（`HAIR_CORS_ORIGINS` 未設定時）
- `uchi.NewRouter()` ラッパーの方針決定（維持するか、直接 `chi.NewRouter()` に戻すか）

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
