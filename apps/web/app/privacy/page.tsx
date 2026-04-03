import Link from "next/link";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "プライバシーポリシー | HairHistory",
  description: "HairHistory のプライバシーに関する方針",
};

export default function PrivacyPage() {
  return (
    <main className="legal-page">
      <p className="legal-page__back">
        <Link href="/">← トップへ戻る</Link>
      </p>
      <h1 className="legal-page__title">プライバシーポリシー</h1>
      <p className="legal-page__meta">最終更新: 2026年4月3日</p>
      <div className="legal-page__body">
        <p>
          本ページは HairHistory（以下「本サービス」）における個人情報および利用に関する方針を説明します。内容は開発段階のものであり、運用開始前に改定される場合があります。
        </p>
        <h2>1. 収集する情報</h2>
        <p>
          本サービスは、匿名のユーザー識別子（URLに含まれるID）、施術履歴の入力内容（日付・施術カテゴリ・サロン名・メモ等）を保存する場合があります。ログインを導入するまでは、メールアドレス等の本人特定情報の取得は行いません（将来変更の可能性があります）。
        </p>
        <h2>2. 利用目的</h2>
        <p>履歴の表示・保存、サービス品質の改善のために利用します。</p>
        <h2>3. 第三者提供</h2>
        <p>
          法令に基づく場合を除き、ご本人の同意なく第三者に個人を特定できる情報を提供しません。ホスティング等のインフラ事業者への委託は行う場合があります。
        </p>
      </div>
    </main>
  );
}
