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
- `main.go` の CORS デフォルト設定に警告ログ追加（`HAIR_CORS_ORIGINS` 未設定時）

## レビューメモ（記録のみ・依頼時点ではコード対応なし）

- **[nits]** `controller/hair_history_create.go` の `json.NewDecoder(r.Body).Decode(in)` 失敗時、エラー文字列がそのままクライアントに返る。既存挙動と同じで今回は OK。将来、クライアント向けにメッセージをラップするか検討してよい。
- **[Q]** `CreateHistory` で URL から `UserID` をセットしたあと `Decode(in)` するが、body に `UserID` 相当の JSON キーがあれば上書きされないか。**struct タグ**: `UserID` は `json:"-"` のため **JSON からは設定されず上書きされない**（`UpdateHistory.HistoryID` も `json:"-"` で同様）。別名キーで別フィールドにマップされるタグが付いていない限り問題なし。
- **[nits]** `usecase/request` の `NewCreateUser()` は現状常に `nil` error を返す。将来バリデーションを足す前提なら問題なし。
- **[Q]** `NewCreateUser()` が引数なしでエラーも常に `nil` だが、戻り値に `error` を残しているのは将来拡張（バリデーション等）を見越してか。**答え（記録）**: その意図であれば一貫してよい。不要ならシグネチャ簡略化も検討の余地あり。
- **[nits]** `privacy/page.tsx`（および同様の法務ページ）の「最終更新」日付がハードコード。運用で更新忘れにだけ注意。現段階では許容範囲。
- **[nits]** `page.tsx` フッターから CONTACT リンクは削除済み。意図的と理解。将来的に何らかの連絡手段を置くかは別途検討してよい。

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
