# HairHistory 実装タスク

このファイルは「進捗 + 次に直すこと」をまとめた1枚です。

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

1. 作成系 API を `201 Created` にする
2. 本番向けに 5xx エラー文言を固定化する
3. `main.go` の起動方式を一本化（`signal.NotifyContext`）

## 優先度中（その次）

- パーマ要件と実装の整合（`ServiceType` の扱いを決定）
- 施術の複数選択UI
- ルート厳密化（将来衝突回避）
- `usecase/request` の `*http.Request` 依存整理
- CONTACT / PRIVACY / TERMS 実リンク化

## 拡張（任意）

- 認証（ログイン）
- 分析・予約・設定
- ダメージサマリー
- 通知 / リマインド
- E2Eテスト整備

---

更新ルール:

- 完了した項目はチェックを付ける
- 仕様変更は `要件定義` / `基本設計` へ反映する
