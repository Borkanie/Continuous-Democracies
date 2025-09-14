import { useQuery } from '@tanstack/react-query';
import { RoundCard } from '../round-card/RoundCard';
import { Search } from '../search/Search';
import styles from './RoundsList.module.css';
import { getAllRounds } from '../../utils/api/rounds';

const { Div, header } = styles;

export const RoundsList = () => {
  const { data } = useQuery({
    queryKey: ['rounds'],
    queryFn: getAllRounds,
  });

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
      {data.map((round) => (
        <RoundCard key={round.voteId} round={round} />
      ))}
    </div>
  );
};
