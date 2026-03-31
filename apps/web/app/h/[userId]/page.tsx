"use client";

import { useParams, useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

import { readApiErrorMessage, readFetchFailureMessage } from "../../../lib/apiError";
import { SharePageQR } from "./SharePageQR";

const api = () => process.env.NEXT_PUBLIC_API_URL ?? "http://127.0.0.1:8080";

const VISIBLE_N = 5;

/** API と一致するコード。表示ラベルはモックに合わせる */
const CATEGORIES: { code: string; label: string }[] = [
  { code: "color", label: "カラー" },
  { code: "bleach", label: "ブリーチ" },
  { code: "straight_perm", label: "縮毛矯正" },
  { code: "other", label: "カット" },
  { code: "treatment", label: "オラプレックス" },
];

function dayToISO(d: string) {
  return `${d}T00:00:00.000Z`;
}

function fmtDateJa(iso: string) {
  const d = iso.slice(0, 10);
  try {
    return new Date(`${d}T12:00:00`).toLocaleDateString("ja-JP", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  } catch {
    return d;
  }
}

type Row = {
  id: string;
  date: string;
  services: string[];
  salonName?: string;
  stylistName?: string;
  memo: string;
};

function labelForService(code: string) {
  return CATEGORIES.find((c) => c.code === code)?.label ?? code;
}

function memoBody(text: string) {
  const t = text.trim();
  if (!t) return "（メモなし）";
  return t;
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
  const [showAll, setShowAll] = useState(false);
  const [copyOk, setCopyOk] = useState(false);
  const [copyErr, setCopyErr] = useState(false);
  const [recordsLoading, setRecordsLoading] = useState(true);
  const [recordsError, setRecordsError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string | null>(null);
  const [savingAdd, setSavingAdd] = useState(false);
  const [savingEdit, setSavingEdit] = useState(false);
  const [savingDeleteId, setSavingDeleteId] = useState<string | null>(null);

  const load = useCallback(
    async (opts?: { showLoading?: boolean }) => {
      const showLoading = opts?.showLoading !== false;
      if (!uid) return;
      if (showLoading) {
        setRecordsLoading(true);
        setRecordsError(null);
      }
      try {
        const res = await fetch(`${api()}/api/users/${uid}/histories`);
        if (res.status === 404) {
          router.push("/");
          return;
        }
        if (!res.ok) {
          const msg = await readApiErrorMessage(res);
          if (showLoading) {
            setRecordsError(msg);
            setList([]);
          } else {
            setNotice(msg);
          }
          return;
        }
        const j = (await res.json()) as { list?: Row[] };
        setList(j.list ?? []);
        if (showLoading) setRecordsError(null);
      } catch (e) {
        const msg = readFetchFailureMessage(e);
        if (showLoading) {
          setRecordsError(msg);
          setList([]);
        } else {
          setNotice(msg);
        }
      } finally {
        if (showLoading) setRecordsLoading(false);
      }
    },
    [uid, router],
  );

  useEffect(() => {
    void load();
  }, [load]);

  async function add(e: React.FormEvent) {
    e.preventDefault();
    if (!uid || savingAdd) return;
    setNotice(null);
    setSavingAdd(true);
    try {
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
      if (!res.ok) {
        setNotice(await readApiErrorMessage(res));
        return;
      }
      setMemo("");
      setSalon("");
      setStylist("");
      setAddDay(new Date().toISOString().slice(0, 10));
      await load({ showLoading: false });
    } catch (e) {
      setNotice(readFetchFailureMessage(e));
    } finally {
      setSavingAdd(false);
    }
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
    if (!editId || savingEdit) return;
    setNotice(null);
    setSavingEdit(true);
    try {
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
      if (!res.ok) {
        setNotice(await readApiErrorMessage(res));
        return;
      }
      setEditId(null);
      await load({ showLoading: false });
    } catch (e) {
      setNotice(readFetchFailureMessage(e));
    } finally {
      setSavingEdit(false);
    }
  }

  async function remove(id: string) {
    if (!globalThis.confirm("この履歴を削除しますか？")) return;
    setNotice(null);
    setSavingDeleteId(id);
    try {
      const res = await fetch(`${api()}/api/histories/${id}`, { method: "DELETE" });
      if (!res.ok) {
        setNotice(await readApiErrorMessage(res));
        return;
      }
      if (editId === id) setEditId(null);
      await load({ showLoading: false });
    } catch (e) {
      setNotice(readFetchFailureMessage(e));
    } finally {
      setSavingDeleteId(null);
    }
  }

  async function copyShareUrl() {
    setCopyErr(false);
    try {
      await navigator.clipboard.writeText(globalThis.location?.href ?? "");
      setCopyOk(true);
      setTimeout(() => setCopyOk(false), 2000);
    } catch {
      setCopyErr(true);
    }
  }

  const visibleList = showAll ? list : list.slice(0, VISIBLE_N);
  const hasMore = list.length > VISIBLE_N;

  function chipRow(value: string, onChange: (c: string) => void) {
    return (
      <div className="archive-chips" role="group" aria-label="カテゴリー">
        {CATEGORIES.map(({ code, label }) => (
          <button
            key={code}
            type="button"
            className={`archive-chip${value === code ? " archive-chip--active" : ""}`}
            onClick={() => onChange(code)}
          >
            {label}
          </button>
        ))}
      </div>
    );
  }

  function editForm() {
    return (
      <form onSubmit={saveEdit} className="archive-edit-form">
        <div className="archive-form-field">
          <span className="archive-form-label">日付</span>
          <div className="archive-date-wrap">
            <input type="date" value={editDay} onChange={(e) => setEditDay(e.target.value)} required />
          </div>
        </div>
        <div className="archive-form-field">
          <span className="archive-form-label">カテゴリー</span>
          {chipRow(editSvc, setEditSvc)}
        </div>
        <div className="archive-form-row2">
          <input value={editSalon} onChange={(e) => setEditSalon(e.target.value)} placeholder="サロン名" aria-label="サロン名" />
          <input
            value={editStylist}
            onChange={(e) => setEditStylist(e.target.value)}
            placeholder="担当スタイリスト"
            aria-label="担当スタイリスト"
          />
        </div>
        <div className="archive-form-field">
          <span className="archive-form-label">メモ</span>
          <textarea
            className="archive-textarea"
            value={editMemo}
            onChange={(e) => setEditMemo(e.target.value)}
            rows={4}
            placeholder="処方や毛髪の状態など"
          />
        </div>
        <div className="archive-record-actions">
          <button type="submit" disabled={savingEdit}>
            {savingEdit ? "保存中…" : "保存"}
          </button>
          <button type="button" disabled={savingEdit} onClick={() => setEditId(null)}>
            キャンセル
          </button>
        </div>
      </form>
    );
  }

  return (
    <main className="archive-main">
      <div className="archive-inner">
        <header className="archive-pagehead">
          <h1>アーカイブ</h1>
          <p className="archive-pagesub">
            カラー、ブリーチ、パーマなど、すべての施術履歴を記録。あなたの髪の履歴を、プロの精度で管理します。
          </p>
        </header>

        {notice ? (
          <div className="app-notice" role="alert">
            <span className="app-notice__text">{notice}</span>
            <button type="button" className="app-notice__dismiss" onClick={() => setNotice(null)}>
              閉じる
            </button>
          </div>
        ) : null}

        <div className="archive-layout">
          <aside className="archive-sidebar">
            <section id="archive-add" className="archive-card">
              <h2>新規施術を記録</h2>
              <form onSubmit={add}>
                <div className="archive-form-field">
                  <label className="archive-form-label" htmlFor="archive-date">
                    日付
                  </label>
                  <div className="archive-date-wrap">
                    <input
                      id="archive-date"
                      type="date"
                      value={addDay}
                      onChange={(e) => setAddDay(e.target.value)}
                      required
                    />
                  </div>
                </div>
                <div className="archive-form-field">
                  <span className="archive-form-label">カテゴリー</span>
                  {chipRow(svc, setSvc)}
                </div>
                <div className="archive-form-row2">
                  <input value={salon} onChange={(e) => setSalon(e.target.value)} placeholder="サロン名" aria-label="サロン名" />
                  <input
                    value={stylist}
                    onChange={(e) => setStylist(e.target.value)}
                    placeholder="担当スタイリスト"
                    aria-label="担当スタイリスト"
                  />
                </div>
                <div className="archive-form-field">
                  <label className="archive-form-label" htmlFor="archive-memo">
                    メモ
                  </label>
                  <textarea
                    id="archive-memo"
                    className="archive-textarea"
                    value={memo}
                    onChange={(e) => setMemo(e.target.value)}
                    rows={5}
                    placeholder="レシピのメモや髪の状態について…"
                  />
                </div>
                <button type="submit" className="archive-btn-primary" disabled={savingAdd}>
                  {savingAdd ? "送信中…" : "記録を追加"}
                </button>
              </form>
            </section>

            <section id="archive-share" className="archive-share">
              <div className="archive-share-head">
                <svg className="archive-share-ico" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" aria-hidden>
                  <path d="M3 3h7v7H3V3zm11 0h7v7h-7V3zM3 14h7v7H3v-7zm14 0h3v3h-3v-3zm-4 4h3v3h-3v-3zm4-4h3v3h-3v-3z" />
                </svg>
                <h3>スタイリストへの共有</h3>
              </div>
              <p className="archive-share-desc">このページのQRコードを生成します。美容師さんに見せるか、読み取ってもらってください。</p>
              <SharePageQR size={128} />
              <button type="button" className="archive-btn-secondary" onClick={() => void copyShareUrl()}>
                {copyOk ? "URLをコピーしました" : "スタイリストに共有"}
              </button>
              {copyErr ? <p className="landing-copy-error">コピーできませんでした。URLを手動で共有してください。</p> : null}
            </section>
          </aside>

          <section className="archive-records" aria-labelledby="records-heading">
            <div className="archive-records-head">
              <h2 id="records-heading">過去の記録</h2>
              <span className="archive-records-count">{recordsLoading ? "読み込み中…" : `${list.length}件の記録が見つかりました`}</span>
            </div>

            {recordsLoading ? (
              <p className="archive-records-loading">履歴を読み込んでいます…</p>
            ) : recordsError ? (
              <div className="archive-records-error">
                <p>{recordsError}</p>
                <button type="button" className="archive-records-retry" onClick={() => void load()}>
                  再読み込み
                </button>
              </div>
            ) : list.length === 0 ? (
              <p className="archive-empty">まだ記録がありません。左のフォームから追加してください。</p>
            ) : (
              <>
                <div className="archive-records-list">
                  {visibleList.map((h) => (
                    <article key={h.id} className="archive-record">
                      {editId === h.id ? (
                        editForm()
                      ) : (
                        <>
                          <div className="archive-record-head">
                            <p className="archive-record-date">{fmtDateJa(h.date)}</p>
                            <span className="archive-record-salon">{h.salonName || "サロン未入力"}</span>
                          </div>
                          <div className="archive-record-tags">
                            {h.services.map((code, i) => (
                              <span key={`${h.id}-${code}-${i}`} className={`archive-tag${i === 0 ? " archive-tag--dark" : ""}`}>
                                {labelForService(code)}
                              </span>
                            ))}
                          </div>
                          <div className="archive-record-body">
                            <div className="archive-record-thumb" aria-hidden />
                            <dl className="archive-record-details">
                              <dt>担当スタイリスト</dt>
                              <dd>{h.stylistName || "—"}</dd>
                              <dt>施術詳細</dt>
                              <dd>{memoBody(h.memo)}</dd>
                            </dl>
                          </div>
                          <div className="archive-record-actions">
                            <button type="button" onClick={() => startEdit(h)}>
                              編集
                            </button>
                            <button type="button" disabled={savingDeleteId !== null} onClick={() => void remove(h.id)}>
                              {savingDeleteId === h.id ? "削除中…" : "削除"}
                            </button>
                          </div>
                        </>
                      )}
                    </article>
                  ))}
                </div>
                {hasMore && !showAll ? (
                  <button type="button" className="archive-loadmore" onClick={() => setShowAll(true)}>
                    すべての履歴を読み込む
                  </button>
                ) : null}
              </>
            )}
          </section>
        </div>
      </div>
    </main>
  );
}
