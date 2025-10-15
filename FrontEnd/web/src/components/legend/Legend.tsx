import type { Slice } from '../chart/PieChart';
import { LegendCard } from '../legend-card/LegendCard';
import styles from './Legend.module.css';

const { Div, bold, cardsContainer } = styles;

type Props = {
  slices: Slice[];
};

export const Legend = (props: Props) => {
  const { slices } = props;

  return (
    <div className={Div}>
      <p className={bold}>Legenda voturilor</p>
      <div className={cardsContainer}>
        {slices.map((slice) => (
          <LegendCard
            key={slice.id}
            text={slice.label}
            color='#4CAF50'
            count={slice.value.count}
            percentage={slice.value.percentage}
          />
        ))}
        {/* <LegendCard
          text='Pentru'
          color='#4CAF50'
          count={23333}
          percentage={33}
        />
        <LegendCard
          text='Impotriva'
          color='#F44336'
          count={233}
          percentage={33}
        />
        <LegendCard
          text='Abtinere'
          color='#FF9800'
          count={233}
          percentage={33}
        />
        <LegendCard
          text='Nespecificat'
          color='#9E9E9E'
          count={233}
          percentage={33}
        />
        <LegendCard
          text='Pentru'
          color='#4CAF50'
          count={23333}
          percentage={33}
        />
        <LegendCard
          text='Impotriva'
          color='#F44336'
          count={233}
          percentage={33}
        />
        <LegendCard
          text='Abtinere'
          color='#FF9800'
          count={233}
          percentage={33}
        />
        <LegendCard
          text='Nespecificat'
          color='#9E9E9E'
          count={233}
          percentage={33}
        />
        <LegendCard
          text='Pentru'
          color='#4CAF50'
          count={23333}
          percentage={33}
        />
        <LegendCard
          text='Impotriva'
          color='#F44336'
          count={233}
          percentage={33}
        />
        <LegendCard
          text='Abtinere'
          color='#FF9800'
          count={233}
          percentage={33}
        />
        <LegendCard
          text='Nespecificat'
          color='#9E9E9E'
          count={233}
          percentage={33}
        /> */}
      </div>
    </div>
  );
};
