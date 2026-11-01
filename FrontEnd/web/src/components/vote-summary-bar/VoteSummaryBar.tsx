import type { MajorityType, Chamber, VoteValue } from '../../utils/types';
import type { GroupedVotes } from '../../utils/hooks/useResultsByRoundId';
import { TOTAL_CHAMBER_MEMBERS } from '../../utils/constants';
import { positionColor, getPassThreshold } from '../../utils/helper';
import { useTheme } from '../../utils/context/ThemeContext';
import styles from './VoteSummaryBar.module.css';

const { wrapper, barArea, counts, count, bar, segment, marker, markerLabel } =
  styles;

const VOTE_VALUES: VoteValue[] = ['Yes', 'No', 'Abstain', 'Absent'];

type Props = {
  grouped: GroupedVotes;
  majorityType: MajorityType;
  chamber: Chamber;
};

export const VoteSummaryBar = (props: Props) => {
  const { grouped, majorityType, chamber } = props;
  const { theme } = useTheme();

  const votesTotal = VOTE_VALUES.reduce(
    (sum, value) => sum + (grouped[value]?.length ?? 0),
    0,
  );
  // The bar spans the whole chamber so the threshold marker stays meaningful
  // even when fewer votes than seats were recorded.
  const barTotal = Math.max(votesTotal, TOTAL_CHAMBER_MEMBERS[chamber]);
  const threshold = getPassThreshold(grouped, majorityType, chamber);
  const markerPosition = Math.min((threshold / barTotal) * 100, 100);

  return (
    <div className={wrapper}>
      <div className={barArea}>
        <div className={counts}>
          {VOTE_VALUES.map((value) => {
            const size = grouped[value]?.length ?? 0;
            return size > 0 ? (
              <span
                key={value}
                className={count}
                style={{
                  width: `${(size / barTotal) * 100}%`,
                  color: positionColor(value, theme),
                }}
              >
                {size}
              </span>
            ) : null;
          })}
        </div>
        <div className={bar}>
          {VOTE_VALUES.map((value) => {
            const size = grouped[value]?.length ?? 0;
            return size > 0 ? (
              <div
                key={value}
                className={segment}
                style={{
                  width: `${(size / barTotal) * 100}%`,
                  backgroundColor: positionColor(value, theme),
                }}
              />
            ) : null;
          })}
          <div className={marker} style={{ left: `${markerPosition}%` }}>
            <span
              className={markerLabel}
              style={markerPosition > 85 ? { transform: 'translateX(-100%)' } : undefined}
            >
              {`Prag: ${threshold}`}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};
