export interface Treatment {
  id: string;
  userId?: string;
  treatedOn: string;
  services: string[];
  salonName?: string | null;
  memo?: string | null;
  cost?: number | null;
  createdAt: string;
  updatedAt: string;
}

export interface TreatmentInput {
  treatedOn: string;
  services: string[];
  salonName?: string;
  memo?: string;
  cost?: number;
}

export type PublicTreatment = Pick<
  Treatment,
  'treatedOn' | 'services' | 'salonName' | 'memo' | 'cost'
>;

export interface ShareLinkRecord {
  id: string;
  token: string;
  expiresAt: string;
  createdAt: string;
}

export interface User {
  id: string;
  name: string;
  email: string;
  createdAt: string;
}
