export interface TreatmentDraft {
  services: string[];
  treatedOn: string;
  salonName: string;
  memo: string;
  cost: string;
}

const draftKey = (id: string) => `hairhistory.treatmentDraft.${id}`;

export const saveTreatmentDraft = (id: string, draft: TreatmentDraft) => {
  try {
    sessionStorage.setItem(draftKey(id), JSON.stringify(draft));
  } catch {
    return;
  }
};

export const readTreatmentDraft = (id: string): TreatmentDraft | null => {
  try {
    const stored = sessionStorage.getItem(draftKey(id));
    if (!stored) return null;
    const parsed = JSON.parse(stored) as Partial<TreatmentDraft>;
    if (!Array.isArray(parsed.services)) return null;
    return {
      services: parsed.services.filter((service) => typeof service === 'string'),
      treatedOn: parsed.treatedOn ?? '',
      salonName: parsed.salonName ?? '',
      memo: parsed.memo ?? '',
      cost: parsed.cost ?? '',
    };
  } catch {
    return null;
  }
};

export const clearTreatmentDraft = (id: string) => {
  try {
    sessionStorage.removeItem(draftKey(id));
  } catch {
    return;
  }
};
