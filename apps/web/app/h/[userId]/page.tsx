"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

const api = () => process.env.NEXT_PUBLIC_API_URL ?? "http://127.0.0.1:8080";

type Row = {
  id: string;
  date: string;
  services: string[];
  salonName?: string;
  stylistName?: string;
  memo: string;
};

const SVC_OPTS: [string, string][] = [
  ["color", "カラー"],
  ["bleach", "ブリーチ"],
  ["straight_perm", "縮毛矯正"],
  ["treatment", "トリートメント"],
  ["other", "その他"],
];

function SvcSelect({
  value,
  onChange,
}: {
  value: string;
  onChange: (v: string) => void;
}) {
  return (
    <select value={value} onChange={(e) => onChange(e.target.value)} aria-label="施術">
      {SVC_OPTS.map(([v, label]) => (
        <option key={v} value={v}>
          {label}
        </option>
      ))}
    </select>
  );
}

export default function HistoryPage() {
  const params = useParams();
  const router = useRouter();
  const uid = typeof params.userId === "string" ? params.userId : params.userId?.[0] ?? "";
  const [list, setList] = useState<Row[]>([]);
  const [svc, setSvc] = useState("color");
  const [memo, setMemo] = useState("");
  const [editId, setEditId] = useState<string | null>(null);
  const [editSvc, setEditSvc] = useState("color");
  const [editMemo, setEditMemo] = useState("");

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

  function startEdit(h: Row) {
    setEditId(h.id);
    setEditSvc(h.services[0] ?? "color");
    setEditMemo(h.memo);
  }

  async function saveEdit(e: React.FormEvent) {
    e.preventDefault();
    if (!editId) return;
    const res = await fetch(`${api()}/api/histories/${editId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ services: [editSvc], memo: editMemo }),
    });
    if (!res.ok) return;
    setEditId(null);
    void load();
  }

  async function remove(id: string) {
    if (!globalThis.confirm("この履歴を削除しますか？")) return;
    const res = await fetch(`${api()}/api/histories/${id}`, { method: "DELETE" });
    if (!res.ok) return;
    if (editId === id) setEditId(null);
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
        <SvcSelect value={svc} onChange={setSvc} />
        <input value={memo} onChange={(e) => setMemo(e.target.value)} placeholder="メモ" />
        <button type="submit">追加</button>
      </form>
      <ul className="hist-list">
        {list.map((h) => (
          <li key={h.id}>
            {editId === h.id ? (
              <form onSubmit={saveEdit} className="hist-edit">
                <SvcSelect value={editSvc} onChange={setEditSvc} />
                <input value={editMemo} onChange={(e) => setEditMemo(e.target.value)} placeholder="メモ" />
                <button type="submit">保存</button>
                <button type="button" onClick={() => setEditId(null)}>
                  キャンセル
                </button>
              </form>
            ) : (
              <>
                <span>
                  {h.date.slice(0, 10)} — {h.services.join(", ")} — {h.memo || "（なし）"}
                </span>
                <span className="hist-actions">
                  <button type="button" onClick={() => startEdit(h)}>
                    編集
                  </button>
                  <button type="button" onClick={() => void remove(h.id)}>
                    削除
                  </button>
                </span>
              </>
            )}
          </li>
        ))}
      </ul>
    </main>
  );
}
