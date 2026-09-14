import { useQuery } from '@tanstack/react-query';
import { RoundCard } from '../round-card/RoundCard';
import { Search } from '../search/Search';
import styles from './RoundsList.module.css';
import { listVotingRounds } from '../../utils/api/rounds';
import type { VotingRound } from '../../utils/types';
import { useNavigate, useParams } from '@tanstack/react-router';
import classNames from 'classnames';
import { Spinner } from '../spinner/Spinner';
import { useState } from 'react';
import { useDebounce } from '../../utils/hooks/useDebounce';
import { EmptyState } from '../empty-state/EmptyState';
import { UiText } from '../ui/text/UiText';
import { useDrawer } from '../../utils/context/DrawerContext';
import { UiButton } from '../ui/button/UiButton';
import { faXmark } from '@fortawesome/free-solid-svg-icons';
import { ScrollableArea } from '../ui/scrollable-area/ScrollableArea';

const { Div, header, inDrawer: inDrawerClass, topSide, headerRight } = styles;

const filterRoundsByKeyword = (
  rounds: VotingRound[] | undefined,
  keyword: string,
): VotingRound[] | undefined => {
  if (!rounds || !keyword) {
    return rounds;
  }

  const lowerKeyword = keyword.toLowerCase();
  return rounds.filter(
    (round) =>
      round.title.toLowerCase().includes(lowerKeyword) ||
      round.description.toLowerCase().includes(lowerKeyword),
  );
};

export const RoundsList = ({ inDrawer = false }) => {
  const [searchTerm, setSearchTerm] = useState('');
  const { debouncedValue, isTyping } = useDebounce(searchTerm);
  const { isDrawerOpen, setIsDrawerOpen } = useDrawer();

  const { data: allRounds, isFetching } = useQuery({
    queryKey: ['rounds'],
    queryFn: () => listVotingRounds(),
  });

  const data = filterRoundsByKeyword(allRounds, debouncedValue);

  const navigate = useNavigate();
  const params = useParams({ strict: false });
  const { roundId } = params;

  const showSpinner = isFetching || isTyping;

  return (
    <div className={classNames(Div, inDrawer && inDrawerClass)}>
      <div className={topSide}>
        <div className={header}>
          <h3>Lista legi</h3>
          <div className={headerRight}>
            <UiText text={'(' + (data?.length || 0) + ')'} />
            {inDrawer && (
              <UiButton
                icon={faXmark}
                onClick={() => setIsDrawerOpen(false)}
                title={'Inchide'}
              />
            )}
          </div>
        </div>
        <Search value={searchTerm} onChange={setSearchTerm} />
      </div>
      <ScrollableArea>
        {showSpinner ? (
          <Spinner />
        ) : (
          <>
            {data && data.length > 0 ? (
              data?.map((round) => (
                <RoundCard
                  key={round.id}
                  round={round}
                  isSelected={roundId === round.id.toString()}
                  onSelect={() => {
                    if (isDrawerOpen) {
                      setIsDrawerOpen(false);
                    }
                    navigate({ to: `round/${round.id}` });
                  }}
                />
              ))
            ) : (
              <EmptyState text={'Nu au fost gasite rezultate.'} />
            )}
          </>
        )}
      </ScrollableArea>
    </div>
  );
};
