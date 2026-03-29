import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = { title: "HairHistory", description: "髪の履歴メモ" };

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
