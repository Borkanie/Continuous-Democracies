import type { Round } from '../../utils/types';
import { DateComponent } from '../date/DateComponent';
import { UiText } from '../ui/text/UiText';
import styles from './RoundCard.module.css';
import classNames from 'classnames';

const { Div, title, footer, header, selected } = styles;

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
    {round.description && (
      <UiText
        text={round.description}
        title={round.description}
        truncate={true}
        maxLines={5}
        size={'medium'}
      />
    )}
    <div className={footer}>
      <DateComponent text={round.voteDate} />
    </div>
  </div>
);
