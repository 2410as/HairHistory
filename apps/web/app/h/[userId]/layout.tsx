import Link from "next/link";

/** `/h/:userId` 用。基本設計のヘッダー（プロジェクト名・説明・トップへの導線） */
export default function HistoryLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="hist-shell">
      <header className="hist-header">
        <p className="hist-brand">
          <Link href="/">HairHistory</Link>
        </p>
        <p className="hist-tagline">美容師さんと一緒にダメージを確認できます。</p>
        <nav>
          <Link href="/">← トップへ</Link>
        </nav>
      </header>
      {children}
    </div>
  );
}
