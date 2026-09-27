import { Box } from '@chakra-ui/react';
import { useEffect, useRef, useState } from 'react';

const GIS_SRC = 'https://accounts.google.com/gsi/client';

// Google の公式ボタンが受け付ける最大幅。
const MAX_WIDTH = 400;

interface CredentialResponse {
  credential?: string;
}

interface GoogleIdentityServices {
  accounts: {
    id: {
      initialize: (config: {
        client_id: string;
        callback: (response: CredentialResponse) => void;
        auto_select?: boolean;
        cancel_on_tap_outside?: boolean;
      }) => void;
      renderButton: (parent: HTMLElement, options: Record<string, unknown>) => void;
      disableAutoSelect: () => void;
    };
  };
}

declare global {
  interface Window {
    google?: GoogleIdentityServices;
  }
}

let loader: Promise<GoogleIdentityServices> | null = null;

const loadGoogleIdentityServices = (): Promise<GoogleIdentityServices> => {
  if (loader) return loader;

  loader = new Promise<GoogleIdentityServices>((resolve, reject) => {
    const settle = () => {
      if (window.google) resolve(window.google);
      else reject(new Error('Google Identity Services did not initialise'));
    };
    const fail = () => reject(new Error(`Failed to load ${GIS_SRC}`));

    const existing = document.querySelector<HTMLScriptElement>(`script[src="${GIS_SRC}"]`);
    if (existing) {
      if (window.google) {
        resolve(window.google);
        return;
      }
      existing.addEventListener('load', settle);
      existing.addEventListener('error', fail);
      return;
    }

    const script = document.createElement('script');
    script.src = GIS_SRC;
    script.async = true;
    script.defer = true;
    script.addEventListener('load', settle);
    script.addEventListener('error', fail);
    document.head.append(script);
  }).catch((error: unknown) => {
    loader = null;
    throw error;
  });

  return loader;
};

interface Props {
  clientId: string;
  onCredential: (idToken: string) => void;
  onError: (message: string) => void;
}

export const GoogleSignInButton = ({ clientId, onCredential, onError }: Props) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const onCredentialRef = useRef(onCredential);
  const onErrorRef = useRef(onError);
  const [width, setWidth] = useState(0);

  useEffect(() => {
    onCredentialRef.current = onCredential;
    onErrorRef.current = onError;
  }, [onCredential, onError]);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    const observer = new ResizeObserver((entries) => {
      const measured = Math.floor(entries[0]?.contentRect.width ?? 0);
      if (measured > 0) setWidth(Math.min(MAX_WIDTH, measured));
    });
    observer.observe(container);
    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    if (width === 0) return;

    let cancelled = false;

    loadGoogleIdentityServices()
      .then((google) => {
        const container = containerRef.current;
        if (cancelled || !container) return;

        google.accounts.id.initialize({
          client_id: clientId,
          auto_select: false,
          cancel_on_tap_outside: true,
          callback: (response) => {
            if (response.credential) onCredentialRef.current(response.credential);
            else onErrorRef.current('Google からトークンを受け取れませんでした。もう一度お試しください。');
          },
        });

        container.replaceChildren();
        google.accounts.id.renderButton(container, {
          type: 'standard',
          theme: 'outline',
          size: 'large',
          shape: 'pill',
          text: 'signin_with',
          logo_alignment: 'center',
          locale: 'ja',
          width,
        });
      })
      .catch(() => {
        if (!cancelled) {
          onErrorRef.current('Google ログインを読み込めませんでした。通信環境を確認してください。');
        }
      });

    return () => {
      cancelled = true;
    };
  }, [clientId, width]);

  return (
    <Box
      ref={containerRef}
      width="100%"
      minHeight="44px"
      display="flex"
      justifyContent="center"
      style={{ colorScheme: 'light' }}
    />
  );
};
