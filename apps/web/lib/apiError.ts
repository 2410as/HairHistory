/** Go API の `render.ErrorJSON` が返す `{ "error": "..." }` を読む */
export async function readApiErrorMessage(res: Response): Promise<string> {
  const text = await res.text();
  if (text) {
    try {
      const j = JSON.parse(text) as { error?: string };
      if (typeof j?.error === "string" && j.error.trim()) {
        return j.error.trim();
      }
    } catch {
      if (text.length < 240) {
        return text.trim();
      }
    }
  }
  return fallbackStatusMessage(res.status);
}

function fallbackStatusMessage(status: number): string {
  const map: Record<number, string> = {
    400: "入力内容を確認してください。",
    401: "認証に失敗しました。",
    403: "この操作は許可されていません。",
    404: "データが見つかりません。",
    500: "サーバーでエラーが発生しました。しばらくしてから再度お試しください。",
    502: "サーバーに接続できませんでした。",
    503: "サービスが一時的に利用できません。",
  };
  return map[status] ?? `通信に失敗しました（${status}）。`;
}

export function readFetchFailureMessage(err: unknown): string {
  if (err instanceof TypeError) {
    return "サーバーに接続できません。API の URL（NEXT_PUBLIC_API_URL）と起動状態を確認してください。";
  }
  return "予期しないエラーが発生しました。";
}
