import type { Position } from './types';

export const VOTERS_TOTAL_NUMBER = 464;

export const SECTION_LABELS: Record<Position, string> = {
  0: 'Da',
  1: 'Nu',
  2: 'Abtinere',
  3: 'Absent',
};

type BreakpointQueries = {
  isMobile: string;
  isTablet: string;
  isDesktop: string;
};

export const defaultQueries: BreakpointQueries = {
  isMobile: '@media only screen and (max-width: 767px)',
  isTablet: '@media only screen and (min-width: 768px) and (max-width: 1023px)',
  isDesktop: '@media only screen and (min-width: 1024px)',
};
