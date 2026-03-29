"use client";

import { useParams, useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

const api = () => process.env.NEXT_PUBLIC_API_URL ?? "http://127.0.0.1:8080";

/** `YYYY-MM-DD` → RFC3339 UTC（API の date 用） */
function dayToISO(d: string) {
  return `${d}T00:00:00.000Z`;
}

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
  const [addDay, setAddDay] = useState(() => new Date().toISOString().slice(0, 10));
  const [svc, setSvc] = useState("color");
  const [salon, setSalon] = useState("");
  const [stylist, setStylist] = useState("");
  const [memo, setMemo] = useState("");
  const [editId, setEditId] = useState<string | null>(null);
  const [editDay, setEditDay] = useState("");
  const [editSvc, setEditSvc] = useState("color");
  const [editSalon, setEditSalon] = useState("");
  const [editStylist, setEditStylist] = useState("");
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
        date: dayToISO(addDay),
        services: [svc],
        salonName: salon,
        stylistName: stylist,
        memo,
      }),
    });
    if (!res.ok) return;
    setMemo("");
    setSalon("");
    setStylist("");
    setAddDay(new Date().toISOString().slice(0, 10));
    void load();
  }

  function startEdit(h: Row) {
    setEditId(h.id);
    setEditDay(h.date.slice(0, 10));
    setEditSvc(h.services[0] ?? "color");
    setEditSalon(h.salonName ?? "");
    setEditStylist(h.stylistName ?? "");
    setEditMemo(h.memo);
  }

  async function saveEdit(e: React.FormEvent) {
    e.preventDefault();
    if (!editId) return;
    const res = await fetch(`${api()}/api/histories/${editId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        date: dayToISO(editDay),
        services: [editSvc],
        salonName: editSalon,
        stylistName: editStylist,
        memo: editMemo,
      }),
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
    <main className="hist-main">
      <h1>マイ履歴</h1>
      <p className="hist-share-hint">このページの URL を美容師さんと共有できます。</p>
      <section className="hist-section" aria-labelledby="add-heading">
        <h2 id="add-heading">新しい履歴を追加</h2>
        <form onSubmit={add} className="hist-add">
        <div className="form-row">
          <label>
            日付{" "}
            <input type="date" value={addDay} onChange={(e) => setAddDay(e.target.value)} required />
          </label>
          <SvcSelect value={svc} onChange={setSvc} />
        </div>
        <div className="form-row">
          <input value={salon} onChange={(e) => setSalon(e.target.value)} placeholder="サロン名" />
          <input value={stylist} onChange={(e) => setStylist(e.target.value)} placeholder="スタイリスト名" />
        </div>
        <div className="form-row">
          <input value={memo} onChange={(e) => setMemo(e.target.value)} placeholder="メモ（ダメージ・薬剤など）" />
          <button type="submit">追加</button>
        </div>
        </form>
      </section>
      <section className="hist-section" aria-labelledby="list-heading">
        <h2 id="list-heading">履歴一覧</h2>
        <ul className="hist-list">
        {list.map((h) => (
          <li key={h.id}>
            {editId === h.id ? (
              <form onSubmit={saveEdit} className="hist-edit">
                <div className="form-row">
                  <label>
                    日付{" "}
                    <input type="date" value={editDay} onChange={(e) => setEditDay(e.target.value)} required />
                  </label>
                  <SvcSelect value={editSvc} onChange={setEditSvc} />
                </div>
                <div className="form-row">
                  <input
                    value={editSalon}
                    onChange={(e) => setEditSalon(e.target.value)}
                    placeholder="サロン名"
                  />
                  <input
                    value={editStylist}
                    onChange={(e) => setEditStylist(e.target.value)}
                    placeholder="スタイリスト名"
                  />
                </div>
                <div className="form-row">
                  <input value={editMemo} onChange={(e) => setEditMemo(e.target.value)} placeholder="メモ" />
                  <button type="submit">保存</button>
                  <button type="button" onClick={() => setEditId(null)}>
                    キャンセル
                  </button>
                </div>
              </form>
            ) : (
              <>
                <span>
                  {h.date.slice(0, 10)} — {h.services.join(", ")} — {h.salonName || "サロン未入力"} —{" "}
                  {h.stylistName || "担当未入力"} — {h.memo || "（メモなし）"}
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
      </section>
    </main>
  );
}
