import { useNavigate } from '@tanstack/react-router';
import type { Slice } from '../chart/PieChart';
import { LegendCard } from '../legend-card/LegendCard';
import styles from './Legend.module.css';
import { UiText } from '../ui/text/UiText';

const { Div, bold, cardsContainer } = styles;

type Props = {
  slices: Slice[];
};

export const Legend = (props: Props) => {
  const { slices } = props;

  const navigate = useNavigate();

  return (
    <div className={Div}>
      <UiText className={bold} text={'Legenda voturilor'} />
      <div className={cardsContainer}>
        {slices.map((slice) => (
          <LegendCard
            key={slice.id}
            text={slice.label}
            color={slice.color}
            count={slice.value.count}
            percentage={slice.value.percentage}
            onClick={() =>
              slice.value.percentage > 0
                ? navigate({ to: `section/${slice.id}` })
                : undefined
            }
          />
        ))}
      </div>
    </div>
  );
};
