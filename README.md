# HairHistoryMemo

モノレポ構成で `web`（Next.js）と `main`（Go API）を分けて開発します。

- `apps/web/`: フロント（Next.js / App Router）
- `apps/main/`: バックエンド（Go / Clean-ish architecture: Controller / Usecase / Domain / Infra）

まずはバックエンドの雛形（層・インタフェース・usecase）を作って、API を通してフロントと繋げていきます。