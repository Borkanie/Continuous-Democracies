import { UiText } from '../ui/text/UiText';
import styles from './LegendCard.module.css';

const { button, marker, countContainer, bold, percentageText, leftSide } =
  styles;

type Props = {
  text: string;
  color: string;
  count: number;
  percentage: number;
  onClick?: () => void;
};

export const LegendCard = (props: Props) => {
  const { text, color, count, percentage, onClick } = props;

  return (
    <button onClick={onClick} className={button}>
      <div className={leftSide}>
        <div className={marker} style={{ backgroundColor: color }}></div>
        <UiText text={text} title={text} />
      </div>
      <div className={countContainer}>
        <p className={bold}>{count}</p>
        <p className={percentageText}>{percentage}%</p>
      </div>
    </button>
  );
};
