import type { Slice } from '../chart/PieChart';
import { LegendCard } from '../legend-card/LegendCard';
import styles from './Legend.module.css';
import { UiText } from '../ui/text/UiText';
import { ScrollableArea } from '../ui/scrollable-area/ScrollableArea';

const { Div, bold, cardsContainer, mh400 } = styles;

type Props = {
  slices: Slice[];
  text: string;
  onSliceClick?: (id: string | number) => void;
};

export const Legend = (props: Props) => {
  const { slices, text, onSliceClick } = props;

  return (
    <div className={Div}>
      <UiText className={bold} text={text} />
      <ScrollableArea className={mh400}>
        <ul className={cardsContainer}>
          {slices.map((slice) => (
            <li key={slice.id}>
              <LegendCard
                text={slice.label}
                color={slice.color}
                count={slice.value.count}
                percentage={slice.value.percentage}
                onClick={
                  slice.value.percentage > 0 && onSliceClick
                    ? () => onSliceClick(slice.id)
                    : undefined
                }
              />
            </li>
          ))}
        </ul>
      </ScrollableArea>
    </div>
  );
};
