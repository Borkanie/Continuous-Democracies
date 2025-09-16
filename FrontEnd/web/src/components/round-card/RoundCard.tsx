import type { Round } from '../../utils/types';
import { DateComp } from '../date/DateComp';
import { Status } from '../status/Status';
import styles from './RoundCard.module.css';
import classNames from 'classnames';

const { Div, title, description, footer, header, selected } = styles;

type Props = {
  round: Round;
  isSelected?: boolean;
  onSelect?: () => void;
};

export const RoundCard = (props: Props) => {
  const { round, isSelected, onSelect } = props;

  // TODO: Use real data from `round` prop when available for date, description, and status
  return (
    <div
      className={classNames(Div, isSelected ? selected : '')}
      onClick={onSelect}
    >
      <div className={header}>
        <h3 className={title}>{round.title}</h3>
      </div>
      <p className={description}>
        Comprehensive legislation to reduce carbon emissions by 50% by 2030
      </p>
      <div className={footer}>
        <DateComp text={round.voteDate} />
        <Status text={'ACTIV'} />
      </div>
    </div>
  );
};
