import Link from "next/link";

/** `/h/:userId` — アーカイブ画面用トップバー */
export default function HistoryLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="archive-shell">
      <header className="archive-topbar">
        <div className="archive-inner archive-topbar-row">
          <Link href="/" className="archive-logo">
            HAIRHISTORY
          </Link>
          <nav className="archive-topnav" aria-label="メイン">
            <span className="archive-topnav-link--active">アーカイブ</span>
            <a href="#archive-add">施術記録</a>
            <span className="archive-topnav-link--muted" title="準備中">
              分析
            </span>
            <span className="archive-topnav-link--muted" title="準備中">
              設定
            </span>
          </nav>
          <div className="archive-top-actions">
            <a href="#archive-share" className="archive-icon-btn" aria-label="共有QRへスクロール">
              <svg className="archive-icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" aria-hidden>
                <path d="M3 3h7v7H3V3zm11 0h7v7h-7V3zM3 14h7v7H3v-7zm14 0h3v3h-3v-3zm-4 4h3v3h-3v-3zm4-4h3v3h-3v-3z" />
              </svg>
            </a>
            <span className="archive-icon-btn" aria-hidden title="アカウント（準備中）">
              <svg className="archive-icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                <circle cx="12" cy="8" r="3.5" />
                <path d="M6 19c0-3.5 3-5 6-5s6 1.5 6 5" />
              </svg>
            </span>
          </div>
        </div>
      </header>
      {children}
    </div>
  );
}
