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
import { useDrawer } from '../../utils/context/DrawerContext';
import { UiButton } from '../ui/button/UiButton';
import { faXmark } from '@fortawesome/free-solid-svg-icons';

const {
  Div,
  header,
  roundsContainer,
  inDrawer: inDrawerClass,
  topSide,
  headerRight,
} = styles;

export const RoundsList = ({ inDrawer = false }) => {
  const [searchTerm, setSearchTerm] = useState('');
  const { debouncedValue, isTyping } = useDebounce(searchTerm);
  const { isDrawerOpen, setIsDrawerOpen } = useDrawer();

  const { data, isFetching } = useQuery({
    queryKey: ['rounds', debouncedValue],
    queryFn: () => getAllRounds(debouncedValue),
  });

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
      <div className={roundsContainer}>
        {showSpinner ? (
          <Spinner />
        ) : (
          <>
            {data && data.length > 0 ? (
              data?.map((round, index) => (
                <RoundCard
                  key={round.voteId + '-' + index}
                  round={round}
                  isSelected={roundId === round.voteId.toString()}
                  onSelect={() => {
                    if (isDrawerOpen) {
                      setIsDrawerOpen(false);
                    }
                    navigate({ to: `round/${round.voteId}` });
                  }}
                />
              ))
            ) : (
              <EmptyState text={'Nu au fost gasite rezultate.'} />
            )}
          </>
        )}
      </div>
    </div>
  );
};
