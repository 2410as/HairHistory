import axios, { AxiosError } from 'axios';

export type ApiErrorCode =
  | 'unauthenticated'
  | 'permission_denied'
  | 'not_found'
  | 'invalid_argument'
  | 'internal'
  | 'network_error';

export class ApiError extends Error {
  readonly code: ApiErrorCode;
  readonly status: number;
  readonly details?: Record<string, string>;

  constructor(
    code: ApiErrorCode,
    message: string,
    status: number,
    details?: Record<string, string>
  ) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.status = status;
    this.details = details;
  }
}

export const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
  headers: { 'Content-Type': 'application/json' },
});

type UnauthenticatedListener = () => void;

const unauthenticatedListeners = new Set<UnauthenticatedListener>();

export const onUnauthenticated = (listener: UnauthenticatedListener) => {
  unauthenticatedListeners.add(listener);
  return () => {
    unauthenticatedListeners.delete(listener);
  };
};

interface ErrorEnvelope {
  error?: {
    code?: string;
    message?: string;
    details?: Record<string, string>;
  };
}

const KNOWN_CODES: ApiErrorCode[] = [
  'unauthenticated',
  'permission_denied',
  'not_found',
  'invalid_argument',
  'internal',
];

const FALLBACK_MESSAGE: Record<ApiErrorCode, string> = {
  unauthenticated: 'ログインが必要です。',
  permission_denied: 'このデータにアクセスする権限がありません。',
  not_found: '対象が見つかりませんでした。',
  invalid_argument: '入力内容を確認してください。',
  internal: 'サーバーでエラーが発生しました。',
  network_error: 'サーバーに接続できませんでした。',
};

const toApiError = (error: AxiosError<ErrorEnvelope>): ApiError => {
  const response = error.response;

  if (!response) {
    return new ApiError('network_error', FALLBACK_MESSAGE.network_error, 0);
  }

  const envelope = response.data?.error;
  const code = KNOWN_CODES.find((known) => known === envelope?.code) ?? 'internal';
  const message = envelope?.message?.trim() || FALLBACK_MESSAGE[code];

  return new ApiError(code, message, response.status, envelope?.details);
};

apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ErrorEnvelope>) => {
    const apiError = toApiError(error);

    if (apiError.code === 'unauthenticated') {
      unauthenticatedListeners.forEach((listener) => listener());
    }

    return Promise.reject(apiError);
  }
);

export const toErrorMessage = (error: unknown): string => {
  if (error instanceof ApiError) return error.message;
  if (error instanceof Error) return error.message;
  return '予期しないエラーが発生しました。';
};
