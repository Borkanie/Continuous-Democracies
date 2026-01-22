import type { Position } from './types';

/* Vote colors (YES, NO, ABSENT, NO VOTE) */
export const positionColor = (position: Position): string => {
  switch (position) {
    case 0:
      return '#36A2EB';
    case 1:
      return '#ff6384';
    case 2:
      return '#ffce56';
    case 3:
      return '#cccccc';
    default:
      return '#cccccc';
  }
};

export const positionLabel = (position: Position): string => {
  switch (position) {
    case 0:
      return 'Da';
    case 1:
      return 'Nu';
    case 2:
      return 'Abtinere';
    case 3:
      return 'Absent';
    default:
      return 'Unknown';
  }
};
