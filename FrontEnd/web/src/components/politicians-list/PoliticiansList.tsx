import type { HydratedVote } from '../../utils/types';
import styles from './PoliticiansList.module.css';
import { ScrollableArea } from '../ui/scrollable-area/ScrollableArea';
import { VoteListItem } from '../vote-list-item/VoteListItem';
import classNames from 'classnames';

const { list, mh60 } = styles;

type Props = {
  vote: HydratedVote[];
  className?: string;
  showPartyIcon?: boolean;
};

export const PoliticiansList = (props: Props) => {
  const { vote, className, showPartyIcon } = props;
  return (
    <ScrollableArea className={classNames(className, mh60)}>
      <ul className={list}>
        {vote?.map((v) => (
          <VoteListItem
            key={v.politicianId}
            vote={v}
            showPartyIcon={showPartyIcon}
          />
        ))}
      </ul>
    </ScrollableArea>
  );
};
