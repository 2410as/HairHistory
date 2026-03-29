"use client";

import { QRCodeSVG } from "qrcode.react";
import { useEffect, useState } from "react";

type Props = { size?: number };

/** 現在のページ URL を QR 化（美容師との共有用） */
export function SharePageQR({ size = 128 }: Props) {
  const [url, setUrl] = useState("");
  useEffect(() => {
    setUrl(globalThis.location?.href ?? "");
  }, []);

  if (!url) {
    return <p className="archive-qr-loading">QR を表示する準備中です…</p>;
  }

  return (
    <>
      <div className="archive-qr-box">
        <QRCodeSVG value={url} size={size} level="M" includeMargin />
      </div>
      <p className="archive-qr-url">
        <code>{url}</code>
      </p>
    </>
  );
}
