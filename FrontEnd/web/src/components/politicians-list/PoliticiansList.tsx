import type { HydratedVote } from '../../utils/types';
import styles from './PoliticiansList.module.css';
import { ScrollableArea } from '../ui/scrollable-area/ScrollableArea';
import classNames from 'classnames';

const { list, card, imageContainer, mh60 } = styles;

type Props = { vote: HydratedVote[]; className?: string };

export const PoliticiansList = (props: Props) => {
  const { vote, className } = props;
  return (
    <ScrollableArea className={classNames(className, mh60)}>
      <ul className={list}>
        {vote?.map(({ politician }) => (
          <li key={politician.id} className={card}>
            <div className={imageContainer}>
              <img
                src={politician.imageUrl || undefined}
                alt={politician.name}
              />
            </div>
            <div>{politician.name}</div>
          </li>
        ))}
      </ul>
    </ScrollableArea>
  );
};
