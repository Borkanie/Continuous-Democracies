import classNames from 'classnames';
import type { VoteResult } from '../../utils/types';
import styles from './PoliticiansList.module.css';

const { list, card, imageContainer } = styles;

type Props = { vote: VoteResult[]; className?: string };

export const PoliticiansList = (props: Props) => {
  const { vote, className } = props;
  return (
    <ul className={classNames(className, list)}>
      {vote?.map(({ politician }) => (
        <li key={politician.id} className={card}>
          <div className={imageContainer}>
            <img src={politician.imageUrl || undefined} alt={politician.name} />
          </div>
          <div>{politician.name}</div>
        </li>
      ))}
    </ul>
  );
};
