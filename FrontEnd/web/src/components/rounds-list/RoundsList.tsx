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
import { EmptyState } from '../empty-state/EmptyState';
import { UiText } from '../ui/text/UiText';

const { Div, header, roundsContainer, separator, top, bottom } = styles;

export const RoundsList = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const { debouncedValue, isTyping } = useDebounce(searchTerm);

  const { data, isFetching } = useQuery({
    queryKey: ['rounds', debouncedValue],
    queryFn: () => getAllRounds(debouncedValue),
    retry: false,
  });

  const navigate = useNavigate();
  const params = useParams({ strict: false });
  const { roundId } = params;

  const showSpinner = isFetching || isTyping;

  return (
    <div className={Div}>
      <div className={header}>
        <h3>Lista legi</h3>
        <UiText text={data?.length || 0} />
      </div>
      <Search value={searchTerm} onChange={setSearchTerm} />
      <div className={roundsContainer}>
        {showSpinner ? (
          <Spinner />
        ) : (
          <>
            <div className={classNames(separator, top)} />
            {data && data.length > 0 ? (
              data?.map((round, index) => (
                <RoundCard
                  key={round.voteId + '-' + index}
                  round={round}
                  isSelected={roundId === round.voteId.toString()}
                  onSelect={() => {
                    navigate({ to: `round/${round.voteId}` });
                  }}
                />
              ))
            ) : (
              <EmptyState text='Nu au fost gasite rezultate.' />
            )}
            <div className={classNames(separator, bottom)} />
          </>
        )}
      </div>
    </div>
  );
};
