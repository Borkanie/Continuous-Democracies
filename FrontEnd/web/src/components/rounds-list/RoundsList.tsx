import { useQuery } from '@tanstack/react-query';
import { RoundCard } from '../round-card/RoundCard';
import { Search } from '../search/Search';
import styles from './RoundsList.module.css';
import { getAllRounds } from '../../utils/api/rounds';
import { useNavigate, useParams } from '@tanstack/react-router';
import classNames from 'classnames';
import { Spinner } from '../spinner/Spinner';
import { useState } from 'react';
import { useDebounce } from '../../utils/hooks/useDebounce';

const { Div, header, roundsContainer, separator, top, bottom } = styles;

export const RoundsList = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const debouncedQuery = useDebounce(searchTerm);

  const { data, isFetching } = useQuery({
    queryKey: ['rounds', debouncedQuery],
    queryFn: () => getAllRounds(debouncedQuery),
    retry: false,
  });

  const navigate = useNavigate();
  const params = useParams({ strict: false });
  const { roundId } = params;

  return (
    <div className={Div}>
      <div className={header}>
        <h3>Lista legi</h3>
        <p>{data?.length}</p>
      </div>
      <Search value={searchTerm} onChange={setSearchTerm} />
      <div className={roundsContainer}>
        {isFetching ? (
          <Spinner />
        ) : (
          <>
            <div className={classNames(separator, top)} />
            {data?.map((round) => (
              <RoundCard
                key={round.voteId}
                round={round}
                isSelected={roundId === round.voteId.toString()}
                onSelect={() => {
                  navigate({ to: `round/${round.voteId}` });
                }}
              />
            ))}
            <div className={classNames(separator, bottom)} />
          </>
        )}
      </div>
    </div>
  );
};
