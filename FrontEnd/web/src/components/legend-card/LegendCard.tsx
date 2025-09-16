import styles from './LegendCard.module.css';

const { Div, marker, countContainer, bold, percentageText, leftSide } = styles;

type Props = {
  text: string;
  color: string;
  count: number;
  percentage: number;
};

export const LegendCard = (props: Props) => {
  const { text, color, count, percentage } = props;

  return (
    <div className={Div}>
      <div className={leftSide}>
        <div className={marker} style={{ backgroundColor: color }}></div>
        <p>{text}</p>
      </div>
      <div className={countContainer}>
        <p className={bold}>{count}</p>
        <p className={percentageText}>{percentage}%</p>
      </div>
    </div>
  );
};
