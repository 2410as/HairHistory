"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

const api = () => process.env.NEXT_PUBLIC_API_URL ?? "http://127.0.0.1:8080";
type Row = { id: string; date: string; services: string[]; memo: string };

export default function HistoryPage() {
  const params = useParams();
  const router = useRouter();
  const uid = typeof params.userId === "string" ? params.userId : params.userId?.[0] ?? "";
  const [list, setList] = useState<Row[]>([]);
  const [svc, setSvc] = useState("color");
  const [memo, setMemo] = useState("");

  const load = useCallback(async () => {
    if (!uid) return;
    const res = await fetch(`${api()}/api/users/${uid}/histories`);
    if (res.status === 404) {
      router.push("/");
      return;
    }
    if (!res.ok) return;
    const j = (await res.json()) as { list?: Row[] };
    setList(j.list ?? []);
  }, [uid, router]);

  useEffect(() => {
    void load();
  }, [load]);

  async function add(e: React.FormEvent) {
    e.preventDefault();
    if (!uid) return;
    const res = await fetch(`${api()}/api/users/${uid}/histories`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        date: new Date().toISOString(),
        services: [svc],
        salonName: "",
        stylistName: "",
        memo,
      }),
    });
    if (!res.ok) return;
    setMemo("");
    void load();
  }

  return (
    <main>
      <p>
        <Link href="/">← トップ</Link>
      </p>
      <h1>マイ履歴</h1>
      <p>この URL を美容師さんと共有できます。</p>
      <form onSubmit={add}>
        <select value={svc} onChange={(e) => setSvc(e.target.value)} aria-label="施術">
          <option value="color">カラー</option>
          <option value="bleach">ブリーチ</option>
          <option value="straight_perm">縮毛矯正</option>
          <option value="treatment">トリートメント</option>
          <option value="other">その他</option>
        </select>
        <input value={memo} onChange={(e) => setMemo(e.target.value)} placeholder="メモ" />
        <button type="submit">追加</button>
      </form>
      <ul>
        {list.map((h) => (
          <li key={h.id}>
            {h.date.slice(0, 10)} — {h.services.join(", ")} — {h.memo || "（なし）"}
          </li>
        ))}
      </ul>
    </main>
  );
}
