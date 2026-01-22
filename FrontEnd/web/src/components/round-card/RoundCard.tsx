import type { Round } from '../../utils/types';
import { DateComponent } from '../date/DateComponent';
import { Status } from '../status/Status';
import styles from './RoundCard.module.css';
import classNames from 'classnames';

const { Div, title, description, footer, header, selected } = styles;

type Props = {
  round: Round;
  isSelected?: boolean;
  onSelect?: () => void;
};

export const RoundCard = ({ round, isSelected, onSelect }: Props) => (
  <div
    className={classNames(Div, isSelected ? selected : '')}
    onClick={onSelect}
  >
    <div className={header}>
      <h3 className={title}>{round.title}</h3>
    </div>
    {round.description && <p className={description}>{round.description}</p>}
    <div className={footer}>
      <DateComponent text={round.voteDate} />
      <Status text={'ACTIV'} />
    </div>
  </div>
);
