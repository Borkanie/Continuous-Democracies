import type { Chamber, VoteValue } from './types';

export const SECTION_LABELS: Record<VoteValue, string> = {
  Yes: 'Da',
  No: 'Nu',
  Abstain: 'Abtinere',
  Absent: 'Absent',
};

/* Steady-state seat counts per chamber, used for absolute/qualified majority math. */
export const TOTAL_CHAMBER_MEMBERS: Record<Chamber, number> = {
  parliament: 330,
  senate: 136,
};
