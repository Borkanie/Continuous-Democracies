import type { VoteValue } from './types';

export const SECTION_LABELS: Record<VoteValue, string> = {
  Yes: 'Da',
  No: 'Nu',
  Abstain: 'Abtinere',
  Absent: 'Absent',
};

/* Matches the steady-state parliament size assumed throughout the backend/seed data. */
export const TOTAL_CHAMBER_MEMBERS = 330;
