import type { VoteValue, MajorityType, HydratedVote, Chamber } from './types';
import type { Theme } from './context/ThemeContext';
import { TOTAL_CHAMBER_MEMBERS } from './constants';

/*
 * Vote colors (YES, NO, ABSTAIN, ABSENT). Dark theme keeps the original,
 * fully-saturated colors (they read well against a dark background); light
 * theme uses lighter pigments of the same hues, since the saturated set
 * looks harsh on a white background.
 */
const POSITION_COLORS_DARK: Record<VoteValue, string> = {
  Yes: '#43a047',
  No: '#e53935',
  Abstain: '#9e9e9e',
  Absent: '#4b5563',
};

const POSITION_COLORS_LIGHT: Record<VoteValue, string> = {
  Yes: '#66bb6a',
  No: '#ef5350',
  Abstain: '#bdbdbd',
  Absent: '#424242',
};

export const positionColor = (value: VoteValue, theme: Theme): string => {
  const palette = theme === 'light' ? POSITION_COLORS_LIGHT : POSITION_COLORS_DARK;
  return palette[value] ?? palette.Abstain;
};

/* Romanian display label for the majority mechanism a round was decided under. */
export const majorityTypeLabel = (majorityType: MajorityType): string => {
  switch (majorityType) {
    case 'simple':
      return 'Majoritate simplă';
    case 'absolute':
      return 'Majoritate absolută';
    case 'qualified':
      return 'Majoritate calificată';
    default:
      return '';
  }
};

/*
 * Neutral gray a "haze" blends toward - a mid-gray for dark theme would read
 * as an odd bright patch against a dark background, so each theme gets its
 * own hazing target rather than a single fixed gray.
 */
const HAZE_COLOR_DARK = '#4a4a4a';
const HAZE_COLOR_LIGHT = '#bdbdbd';

/*
 * Blends a 6-digit hex color toward a neutral gray by `amount` (0-1), used to
 * cover unfocused slices in a grayed-out haze while the highlighted one stays
 * at full color. Blending toward gray (rather than toward black, or toward
 * transparent/white) keeps a hint of the underlying hue so slices stay
 * distinguishable, without looking either "darkened" or "faded."
 */
export const haze = (hexColor: string, theme: Theme, amount = 0.65): string => {
  const hazeColor = theme === 'light' ? HAZE_COLOR_LIGHT : HAZE_COLOR_DARK;
  const original = hexColor.replace('#', '');
  const gray = hazeColor.replace('#', '');
  const channel = (start: number) => {
    const originalValue = parseInt(original.substring(start, start + 2), 16);
    const grayValue = parseInt(gray.substring(start, start + 2), 16);
    const blended = Math.round(originalValue * (1 - amount) + grayValue * amount);
    return blended.toString(16).padStart(2, '0');
  };
  return `#${channel(0)}${channel(2)}${channel(4)}`;
};

/*
 * Minimum number of Yes votes needed for a round to pass, based on which of
 * the three Romanian Parliament majority mechanisms applies: simple (more Yes
 * than No), absolute (Yes over half the chamber), or qualified (Yes at least
 * 2/3 of the chamber, used for constitutional revision or presidential
 * impeachment).
 */
export const getPassThreshold = (
  grouped: Record<VoteValue, HydratedVote[]>,
  majorityType: MajorityType,
  chamber: Chamber,
): number => {
  const noCount = grouped.No?.length ?? 0;
  const totalMembers = TOTAL_CHAMBER_MEMBERS[chamber];

  switch (majorityType) {
    case 'absolute':
      return Math.floor(totalMembers / 2) + 1;
    case 'qualified':
      return Math.ceil((totalMembers * 2) / 3);
    case 'simple':
    default:
      return noCount + 1;
  }
};

/* Determines whether a round passed, by comparing Yes votes to the pass threshold. */
export const getVoteStatus = (
  grouped: Record<VoteValue, HydratedVote[]>,
  majorityType: MajorityType,
  chamber: Chamber,
): 'PASSED' | 'FAILED' => {
  const yesCount = grouped.Yes?.length ?? 0;
  return yesCount >= getPassThreshold(grouped, majorityType, chamber)
    ? 'PASSED'
    : 'FAILED';
};

export const positionLabel = (value: VoteValue): string => {
  switch (value) {
    case 'Yes':
      return 'Da';
    case 'No':
      return 'Nu';
    case 'Abstain':
      return 'Abtinere';
    case 'Absent':
      return 'Absent';
    default:
      return 'Unknown';
  }
};
