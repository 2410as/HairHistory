"use client";

import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState, type ReactNode } from "react";

import { readApiErrorMessage, readFetchFailureMessage } from "../lib/apiError";

const api = () => process.env.NEXT_PUBLIC_API_URL ?? "http://127.0.0.1:8080";
const LS = "hairhistoryUserId";

function Svg({ cls, children }: { cls: string; children: ReactNode }) {
  return (
    <svg className={cls} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" aria-hidden>
      {children}
    </svg>
  );
}

export default function HomePage() {
  const router = useRouter();
  const [saved, setSaved] = useState<string | null>(null);
  const [startError, setStartError] = useState<string | null>(null);
  const [starting, setStarting] = useState(false);
  useEffect(() => setSaved(localStorage.getItem(LS)), []);

  async function start() {
    setStartError(null);
    setStarting(true);
    try {
      const res = await fetch(`${api()}/api/users`, { method: "POST" });
      if (!res.ok) {
        setStartError(await readApiErrorMessage(res));
        return;
      }
      const j = (await res.json()) as { ent?: { id: string } };
      const id = j.ent?.id;
      if (!id) {
        setStartError("ユーザーIDを取得できませんでした。");
        return;
      }
      localStorage.setItem(LS, id);
      setSaved(id);
      router.push(`/h/${id}`);
    } catch (e) {
      setStartError(readFetchFailureMessage(e));
    } finally {
      setStarting(false);
    }
  }

  function goQR() {
    if (saved) router.push(`/h/${saved}`);
    else document.getElementById("landing-start")?.scrollIntoView({ behavior: "smooth", block: "center" });
  }

  return (
    <div className="landing-shell">
      {startError ? (
        <div className="app-notice landing-notice-banner" role="alert">
          <span className="app-notice__text">{startError}</span>
          <button type="button" className="app-notice__dismiss" onClick={() => setStartError(null)}>
            閉じる
          </button>
        </div>
      ) : null}
      <header className="landing-header">
        <div className="landing-inner landing-header-inner">
          <Link href="/" className="landing-brand">
            HAIRHISTORY
          </Link>
          <div className="landing-header-end">
            <nav className="landing-nav" aria-label="メイン">
              <Link href="/" className="landing-nav-link landing-nav-link--active">
                ホーム
              </Link>
              {saved ? (
                <Link href={`/h/${saved}`} className="landing-nav-link">
                  履歴
                </Link>
              ) : (
                <a
                  href="#landing-start"
                  className="landing-nav-link landing-nav-link--hint"
                  title="先に「はじめる」でIDを発行すると、ここからマイ履歴（/h/…）を開けます"
                >
                  履歴
                </a>
              )}
              <span className="landing-nav-link landing-nav-link--disabled" title="準備中">
                予約
              </span>
              <span className="landing-nav-link landing-nav-link--disabled" title="準備中">
                設定
              </span>
            </nav>
            <button
              type="button"
              className="landing-qr-btn"
              onClick={goQR}
              aria-label="マイ履歴のQRを表示"
              title={saved ? "マイ履歴へ" : "タップで「はじめる」へスクロール（ID発行後はマイ履歴へ）"}
            >
              <Svg cls="landing-qr-ico">
                <path d="M3 3h7v7H3V3zm11 0h7v7h-7V3zM3 14h7v7H3v-7zm14 0h3v3h-3v-3zm-4 4h3v3h-3v-3zm4-4h3v3h-3v-3z" />
              </Svg>
            </button>
          </div>
        </div>
      </header>

      <section className="landing-hero-band" aria-label="メインビジュアル">
        <div className="landing-inner landing-hero-center">
          <h1 className="landing-hero-title">髪の履歴を、美容師さんと共有しよう。</h1>
          <p className="landing-hero-lead">
            過去の施術、ダメージの記憶。スマホに保存されたあなたの髪の履歴を、QRコード一つでプロの手に。
          </p>
          <div className="landing-hero-cta">
            <button type="button" id="landing-start" className="landing-btn landing-btn--black" onClick={start} disabled={starting}>
              {starting ? "発行中…" : "はじめる（IDを発行する）"}
            </button>
            <a href="#features" className="landing-btn landing-btn--muted">
              使い方を見る
            </a>
          </div>
          {saved ? (
            <p className="landing-hero-saved">
              <Link href={`/h/${saved}`}>保存済みのマイ履歴を開く</Link>
            </p>
          ) : null}
        </div>
      </section>

      <section id="features" className="landing-features">
        <div className="landing-inner landing-features-grid">
          <div className="landing-features-photo">
            <Image
              src="/landing-feature.png"
              alt="髪のラインアート"
              width={969}
              height={1024}
              className="landing-features-img"
              sizes="(max-width: 960px) 100vw, 50vw"
              priority
            />
          </div>
          <ol className="landing-step-list">
            <li>
              <span className="landing-step-num">01</span>
              <div>
                <h2 className="landing-step-title">履歴を記録</h2>
                <p className="landing-step-desc">
                  カラー、パーマ、トリートメント。日々の施術内容やダメージの状態を直感的に記録します。
                </p>
              </div>
            </li>
            <li>
              <span className="landing-step-num">02</span>
              <div>
                <h2 className="landing-step-title">QRを表示</h2>
                <p className="landing-step-desc">
                  生成されたQRコードを美容師さんに提示するだけ。複雑な説明はもう必要ありません。
                </p>
              </div>
            </li>
            <li>
              <span className="landing-step-num">03</span>
              <div>
                <h2 className="landing-step-title">最適な施術へ</h2>
                <p className="landing-step-desc">
                  過去の薬剤履歴に基づき、今のあなたの髪に最適なアプローチをプロが導き出します。
                </p>
              </div>
            </li>
          </ol>
        </div>
      </section>

      <section className="landing-chronicle" aria-labelledby="chronicle-heading">
        <div className="landing-inner">
          <h2 id="chronicle-heading" className="landing-chronicle-en">
            Chronicle of Beauty
          </h2>
          <p className="landing-chronicle-ja">美しさは、時間の積み重ねから生まれる。</p>
          <div className="landing-timeline">
            <div className="landing-tl-row">
              <div className="landing-tl-side">
                <div className="landing-tl-meta">
                  <time dateTime="2023-11-12">2023.11.12</time>
                  <strong>High-Tone Bleach</strong>
                </div>
              </div>
              <div className="landing-tl-side">
                <div className="landing-tl-card">
                  ダメージレベル: 3/5。18レベルまでリフトアップ。補修剤を併用。
                </div>
              </div>
            </div>
            <div className="landing-tl-row landing-tl-row--alt">
              <div className="landing-tl-side">
                <div className="landing-tl-meta">
                  <time dateTime="2024-01-05">2024.01.05</time>
                  <strong>Deep Gray Toning</strong>
                </div>
              </div>
              <div className="landing-tl-side">
                <div className="landing-tl-card">
                  弱酸性カラーを使用。毛先の乾燥をカバーするオイルケアを推奨。
                </div>
              </div>
            </div>
            <div className="landing-tl-row">
              <div className="landing-tl-side">
                <div className="landing-tl-meta">
                  <time dateTime="2024-03-20">2024.03.20</time>
                  <strong>Acidic Treatment</strong>
                </div>
              </div>
              <div className="landing-tl-side">
                <div className="landing-tl-card">
                  酸熱トリートメントで広がりを抑制。髪の芯まで栄養を補給。
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="landing-privacy" id="privacy" aria-labelledby="privacy-heading">
        <div className="landing-inner landing-privacy-box">
          <div className="landing-lock" aria-hidden>
            <Svg cls="landing-lock-svg">
              <rect x="5" y="11" width="14" height="10" rx="1" />
              <path d="M8 11V7a4 4 0 0 1 8 0v4" />
            </Svg>
          </div>
          <h2 id="privacy-heading" className="landing-privacy-title">
            あなたのデータは、あなただけのもの。
          </h2>
          <p className="landing-privacy-text">
            HairHistoryは、あなたのプライバシーを最優先に設計されています。記録されたデータは暗号化され、あなたがQRコードを提示した瞬間だけ、信頼できるプロフェッショナルと共有されます。
          </p>
          <a className="landing-privacy-cta" href="#footer-legal">
            PRIVACY POLICY
          </a>
        </div>
      </section>

      <section className="landing-ready" aria-labelledby="ready-heading">
        <div className="landing-inner landing-ready-inner">
          <p className="landing-ready-label">READY TO BEGIN</p>
          <h2 id="ready-heading" className="landing-ready-title">
            理想のスタイルへの最短距離を。
          </h2>
          <div className="landing-register-card">
            <div className="landing-register-qr" aria-hidden>
              <Svg cls="landing-register-qr-svg">
                <path d="M3 3h7v7H3V3zm11 0h7v7h-7V3zM3 14h7v7H3v-7zm14 0h3v3h-3v-3zm-4 4h3v3h-3v-3zm4-4h3v3h-3v-3z" />
              </Svg>
            </div>
            <div className="landing-register-copy">
              <h3 className="landing-register-h">今すぐHairHistory IDを発行</h3>
              <p className="landing-register-p">登録は無料です。あなたの髪のカルテを今日から作り始めましょう。</p>
              <button type="button" className="landing-btn landing-btn--black landing-btn--wide" onClick={start} disabled={starting}>
                {starting ? "発行中…" : "REGISTER NOW"}
              </button>
            </div>
          </div>
        </div>
      </section>

      <footer className="landing-footer" id="footer-legal">
        <div className="landing-inner landing-footer-inner">
          <span className="landing-footer-brand">HAIRHISTORY</span>
          <nav className="landing-footer-nav" aria-label="フッター">
            <a href="#privacy">PRIVACY</a>
            <a href="#footer-legal">TERMS</a>
            <a href="mailto:support@example.com">CONTACT</a>
          </nav>
          <span className="landing-footer-copy">© 2026 HAIRHISTORY. ALL RIGHTS RESERVED.</span>
        </div>
      </footer>
    </div>
  );
}
