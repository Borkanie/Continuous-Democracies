import { useQuery } from '@tanstack/react-query';
import { RoundCard } from '../round-card/RoundCard';
import { Search } from '../search/Search';
import styles from './RoundsList.module.css';
import { getAllRounds } from '../../utils/api/rounds';
import { useState } from 'react';

const { Div, header, roundsContainer } = styles;

export const RoundsList = () => {
  const { data } = useQuery({
    queryKey: ['rounds'],
    queryFn: getAllRounds,
  });

  const [selectedId, setSelectedId] = useState<number | null>(null);

  if (!data) {
    // TODO: Add empty state here
    return <></>;
  }

  return (
    <div className={Div}>
      <div className={header}>
        <h2>Lista legi</h2>
        <p>244</p>
      </div>
      <Search />
      <div className={roundsContainer}>
        {data.map((round) => (
          <RoundCard
            key={round.voteId}
            round={round}
            isSelected={selectedId === round.voteId}
            onSelect={() => setSelectedId(round.voteId)}
          />
        ))}
      </div>
    </div>
  );
};
