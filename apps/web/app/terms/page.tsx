import Link from "next/link";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "利用規約 | HairHistory",
  description: "HairHistory の利用に関する条件",
};

export default function TermsPage() {
  return (
    <main className="legal-page">
      <p className="legal-page__back">
        <Link href="/">← トップへ戻る</Link>
      </p>
      <h1 className="legal-page__title">利用規約</h1>
      <p className="legal-page__meta">最終更新: 2026年4月3日</p>
      <div className="legal-page__body">
        <p>
          本規約は、HairHistory（以下「本サービス」）の利用条件を定めるものです。本サービスを利用することで、本規約に同意したものとみなします。内容は開発段階のものであり、改定される場合があります。
        </p>
        <h2>1. サービス内容</h2>
        <p>本サービスは、髪の施術履歴を記録・共有するための機能を提供します。機能は予告なく変更・中止される場合があります。</p>
        <h2>2. 禁止事項</h2>
        <p>法令違反、他人の権利侵害、サービス運営を妨げる行為、その他当社が不適切と判断する行為を禁止します。</p>
        <h2>3. 免責</h2>
        <p>
          本サービスは現状有姿で提供されます。医療・美容に関する判断は専門家にご相談ください。本サービスの利用により生じた損害について、当社に故意または重過失がある場合を除き責任を負いません。
        </p>
        <h2>4. 準拠法</h2>
        <p>本規約は日本法を準拠法とします。</p>
      </div>
    </main>
  );
}
