import { useCallback, useEffect, useState } from 'react';
import { apiClient } from '../lib/apiClient';
import type { PublicTreatment, ShareLinkRecord } from '../types';

export const buildShareUrl = (token: string) =>
  `${window.location.origin}/share/${token}`;

interface UseSharesResult {
  shares: ShareLinkRecord[];
  isLoading: boolean;
  error: Error | null;
  isCreating: boolean;
  actionError: Error | null;
  create: () => Promise<ShareLinkRecord | null>;
  revoke: (token: string) => Promise<boolean>;
}

export const useShares = (): UseSharesResult => {
  const [shares, setShares] = useState<ShareLinkRecord[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  const [actionError, setActionError] = useState<Error | null>(null);

  const load = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.get<{ shares: ShareLinkRecord[] }>('/api/shares');
      setShares(response.data.shares ?? []);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  const create = useCallback(async () => {
    setIsCreating(true);
    setActionError(null);
    try {
      const response = await apiClient.post<ShareLinkRecord>('/api/shares', {
        expiresInHours: 168,
      });
      setShares((prev) => [response.data, ...prev]);
      return response.data;
    } catch (err) {
      setActionError(err instanceof Error ? err : new Error('Unknown error'));
      return null;
    } finally {
      setIsCreating(false);
    }
  }, []);

  const revoke = useCallback(async (token: string) => {
    setActionError(null);
    try {
      await apiClient.delete(`/api/shares/${token}`);
      setShares((prev) => prev.filter((share) => share.token !== token));
      return true;
    } catch (err) {
      setActionError(err instanceof Error ? err : new Error('Unknown error'));
      return false;
    }
  }, []);

  return { shares, isLoading, error, isCreating, actionError, create, revoke };
};

interface UsePublicShareResult {
  treatments: PublicTreatment[];
  isLoading: boolean;
  error: Error | null;
}

export const usePublicShare = (token: string | undefined): UsePublicShareResult => {
  const [treatments, setTreatments] = useState<PublicTreatment[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    if (!token) {
      setIsLoading(false);
      return;
    }

    let active = true;
    setIsLoading(true);
    setError(null);

    apiClient
      .get<{ treatments: PublicTreatment[] }>(`/api/public/shares/${token}`)
      .then((response) => {
        if (active) setTreatments(response.data.treatments ?? []);
      })
      .catch((err: unknown) => {
        if (active) setError(err instanceof Error ? err : new Error('Unknown error'));
      })
      .finally(() => {
        if (active) setIsLoading(false);
      });

    return () => {
      active = false;
    };
  }, [token]);

  return { treatments, isLoading, error };
};
