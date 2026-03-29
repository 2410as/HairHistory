"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

const api = () => process.env.NEXT_PUBLIC_API_URL ?? "http://127.0.0.1:8080";
const LS = "hairhistoryUserId";

export default function HomePage() {
  const router = useRouter();
  const [saved, setSaved] = useState<string | null>(null);
  useEffect(() => setSaved(localStorage.getItem(LS)), []);

  async function start() {
    const res = await fetch(`${api()}/api/users`, { method: "POST" });
    if (!res.ok) return;
    const j = (await res.json()) as { ent?: { id: string } };
    const id = j.ent?.id;
    if (!id) return;
    localStorage.setItem(LS, id);
    setSaved(id);
    router.push(`/h/${id}`);
  }

  return (
    <main>
      <h1>HairHistory</h1>
      <p>美容師さんと一緒にダメージを確認できます。</p>
      <p>
        <button type="button" onClick={start}>
          はじめる
        </button>
      </p>
      {saved ? (
        <p>
          <Link href={`/h/${saved}`}>保存済みのマイ履歴へ</Link>
        </p>
      ) : null}
    </main>
  );
}
