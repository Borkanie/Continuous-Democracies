import { useQuery } from '@tanstack/react-query';
import { RoundCard } from '../round-card/RoundCard';
import { Search } from '../search/Search';
import styles from './RoundsList.module.css';
import { getAllRounds } from '../../utils/api/rounds';
import { useNavigate, useParams } from '@tanstack/react-router';
import classNames from 'classnames';

const { Div, header, roundsContainer, separator, top, bottom } = styles;

export const RoundsList = () => {
  const { data } = useQuery({
    queryKey: ['rounds'],
    queryFn: getAllRounds,
  });

  const navigate = useNavigate();
  const params = useParams({ strict: false });
  const { roundId } = params;

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
        <div className={classNames(separator, top)}></div>
        {data.map((round) => (
          <RoundCard
            key={round.voteId}
            round={round}
            isSelected={roundId === round.voteId.toString()}
            onSelect={() => {
              navigate({ to: `round/${round.voteId}` });
            }}
          />
        ))}
        <div className={classNames(separator, bottom)}></div>
      </div>
    </div>
  );
};
