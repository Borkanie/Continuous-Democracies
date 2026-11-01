import { useMemo } from 'react';
import type { HydratedVote, Party, VoteValue } from '../../utils/types';
import { PartyIcon } from '../ui/party-icon/PartyIcon';
import { positionColor } from '../../utils/helper';
import { useTheme } from '../../utils/context/ThemeContext';
import styles from './PartyVoteRows.module.css';

const { list, row, yesCell, centerCell, restCell, barSegment, countLabel, bar } =
  styles;

const RIGHT_VALUES: VoteValue[] = ['No', 'Abstain', 'Absent'];

type PartyCounts = {
  party: Party;
  size: number;
  counts: Record<VoteValue, number>;
};

type Props = {
  votes: HydratedVote[];
  onPartyClick?: (partyId: string) => void;
};

export const PartyVoteRows = (props: Props) => {
  const { votes, onPartyClick } = props;
  const { theme } = useTheme();

  const parties = useMemo(() => groupByParty(votes), [votes]);
  const maxSize = parties[0]?.size ?? 1;

  const renderSegment = (item: PartyCounts, value: VoteValue) => {
    const size = item.counts[value];
    if (size === 0) {
      return null;
    }
    const color = positionColor(value, theme);
    return (
      <div
        key={value}
        className={barSegment}
        style={{ width: `${(size / maxSize) * 100}%` }}
      >
        <span className={countLabel} style={{ color }}>
          {size}
        </span>
        <div className={bar} style={{ backgroundColor: color }} />
      </div>
    );
  };

  return (
    <div className={list}>
      {parties.map((item) => (
        <button
          key={item.party.id}
          type={'button'}
          className={row}
          title={item.party.name}
          onClick={() => onPartyClick?.(item.party.id)}
        >
          <div className={yesCell}>{renderSegment(item, 'Yes')}</div>
          <div className={centerCell}>
            <PartyIcon party={item.party} size={36} />
          </div>
          <div className={restCell}>
            {RIGHT_VALUES.map((value) => renderSegment(item, value))}
          </div>
        </button>
      ))}
    </div>
  );
};

/* Counts each vote value per party, largest party first. Votes without a party are skipped. */
const groupByParty = (votes: HydratedVote[]): PartyCounts[] => {
  const byId = new Map<string, PartyCounts>();
  votes.forEach((vote) => {
    if (!vote.party) {
      return;
    }
    let entry = byId.get(vote.party.id);
    if (!entry) {
      entry = {
        party: vote.party,
        size: 0,
        counts: { Yes: 0, No: 0, Abstain: 0, Absent: 0 },
      };
      byId.set(vote.party.id, entry);
    }
    entry.size += 1;
    entry.counts[vote.value] += 1;
  });
  return Array.from(byId.values()).sort((a, b) => b.size - a.size);
};
