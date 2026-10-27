import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from 'react';
import { ApiError, apiClient, onUnauthenticated } from '../lib/apiClient';
import type { User } from '../types';

interface AuthContextValue {
  user: User | null;
  isLoading: boolean;
  isSessionExpired: boolean;
  devLogin: (email: string, name: string) => Promise<void>;
  googleLogin: (idToken: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSessionExpired, setIsSessionExpired] = useState(false);
  const hasLiveSessionRef = useRef(false);

  useEffect(() => {
    hasLiveSessionRef.current = user !== null;
  }, [user]);

  useEffect(
    () =>
      onUnauthenticated(() => {
        if (!hasLiveSessionRef.current) return;
        hasLiveSessionRef.current = false;
        setIsSessionExpired(true);
        setUser(null);
      }),
    []
  );

  useEffect(() => {
    let active = true;

    apiClient
      .get<User>('/api/auth/me')
      .then((response) => {
        if (active) setUser(response.data);
      })
      .catch((error: unknown) => {
        if (!active) return;
        setUser(null);
        if (!(error instanceof ApiError) || error.code !== 'unauthenticated') {
          console.error(error);
        }
      })
      .finally(() => {
        if (active) setIsLoading(false);
      });

    return () => {
      active = false;
    };
  }, []);

  const devLogin = useCallback(async (email: string, name: string) => {
    const response = await apiClient.post<User>('/api/auth/dev-login', { email, name });
    setIsSessionExpired(false);
    setUser(response.data);
  }, []);

  const googleLogin = useCallback(async (idToken: string) => {
    const response = await apiClient.post<User>('/api/auth/google', { idToken });
    setIsSessionExpired(false);
    setUser(response.data);
  }, []);

  const logout = useCallback(async () => {
    hasLiveSessionRef.current = false;
    await apiClient.post('/api/auth/logout');
    setIsSessionExpired(false);
    setUser(null);
  }, []);

  const value = useMemo(
    () => ({ user, isLoading, isSessionExpired, devLogin, googleLogin, logout }),
    [user, isLoading, isSessionExpired, devLogin, googleLogin, logout]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = (): AuthContextValue => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used inside AuthProvider');
  }
  return context;
};
