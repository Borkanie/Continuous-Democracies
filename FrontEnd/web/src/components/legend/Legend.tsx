import { useNavigate } from '@tanstack/react-router';
import type { Slice } from '../chart/PieChart';
import { LegendCard } from '../legend-card/LegendCard';
import styles from './Legend.module.css';

const { Div, bold, cardsContainer } = styles;

type Props = {
  slices: Slice[];
};

export const Legend = (props: Props) => {
  const { slices } = props;

  const navigate = useNavigate();

  return (
    <div className={Div}>
      <p className={bold}>Legenda voturilor</p>
      <div className={cardsContainer}>
        {slices.map((slice) => (
          <LegendCard
            key={slice.id}
            text={slice.label}
            color={slice.color}
            count={slice.value.count}
            percentage={slice.value.percentage}
            onClick={() => navigate({ to: `section/${slice.id}` })}
          />
        ))}
      </div>
    </div>
  );
};
