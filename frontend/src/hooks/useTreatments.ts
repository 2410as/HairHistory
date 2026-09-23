import { useCallback, useEffect, useState } from 'react';
import { apiClient } from '../lib/apiClient';
import type { Treatment, TreatmentInput } from '../types';

interface UseTreatmentsResult {
  data: Treatment[];
  isLoading: boolean;
  error: Error | null;
  refetch: () => Promise<void>;
}

export const useTreatments = (): UseTreatmentsResult => {
  const [data, setData] = useState<Treatment[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const load = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.get<{ treatments: Treatment[] }>('/api/treatments');
      setData(response.data.treatments ?? []);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  return { data, isLoading, error, refetch: load };
};

interface UseTreatmentResult {
  data: Treatment | null;
  isLoading: boolean;
  error: Error | null;
}

export const useTreatment = (id: string | undefined): UseTreatmentResult => {
  const [data, setData] = useState<Treatment | null>(null);
  const [isLoading, setIsLoading] = useState(Boolean(id));
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    if (!id) {
      setData(null);
      setError(null);
      setIsLoading(false);
      return;
    }

    let active = true;
    setIsLoading(true);
    setError(null);

    apiClient
      .get<Treatment>(`/api/treatments/${id}`)
      .then((response) => {
        if (active) setData(response.data);
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
  }, [id]);

  return { data, isLoading, error };
};

export const createTreatment = async (input: TreatmentInput): Promise<Treatment> => {
  const response = await apiClient.post<Treatment>('/api/treatments', input);
  return response.data;
};

export const updateTreatment = async (
  id: string,
  input: TreatmentInput
): Promise<Treatment> => {
  const response = await apiClient.put<Treatment>(`/api/treatments/${id}`, input);
  return response.data;
};

export const deleteTreatment = async (id: string): Promise<true> => {
  await apiClient.delete(`/api/treatments/${id}`);
  return true;
};

interface MutationState {
  isSubmitting: boolean;
  error: Error | null;
}

const useMutation = <TArgs extends unknown[], TResult>(
  mutate: (...args: TArgs) => Promise<TResult>
) => {
  const [state, setState] = useState<MutationState>({ isSubmitting: false, error: null });

  const run = useCallback(
    async (...args: TArgs): Promise<TResult | null> => {
      setState({ isSubmitting: true, error: null });
      try {
        const result = await mutate(...args);
        setState({ isSubmitting: false, error: null });
        return result;
      } catch (err) {
        const normalized = err instanceof Error ? err : new Error('Unknown error');
        setState({ isSubmitting: false, error: normalized });
        return null;
      }
    },
    [mutate]
  );

  return { run, ...state };
};

export const useCreateTreatment = () => useMutation(createTreatment);
export const useUpdateTreatment = () => useMutation(updateTreatment);
export const useDeleteTreatment = () => useMutation(deleteTreatment);
