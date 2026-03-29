"use client";

import { QRCodeSVG } from "qrcode.react";
import { useEffect, useState } from "react";

/** 現在のページ URL を QR 化（美容師との共有用） */
export function SharePageQR() {
  const [url, setUrl] = useState("");
  useEffect(() => {
    setUrl(globalThis.location?.href ?? "");
  }, []);

  if (!url) {
    return <p className="hist-qr-loading">QR を表示する準備中です…</p>;
  }

  return (
    <section className="hist-section hist-qr" aria-labelledby="qr-heading">
      <h2 id="qr-heading">共有用 QR コード</h2>
      <p className="hist-qr-hint">美容師さんにこの画面を見せるか、スマホで読み取ってもらってください。</p>
      <div className="hist-qr-box">
        <QRCodeSVG value={url} size={168} level="M" includeMargin />
      </div>
      <p className="hist-qr-url">
        <code>{url}</code>
      </p>
    </section>
  );
}
