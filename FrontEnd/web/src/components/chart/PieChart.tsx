import { useRef } from 'react';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Pie } from 'react-chartjs-2';
import ChartDataLabels from 'chartjs-plugin-datalabels';
import styles from './PieChart.module.css';

const { pie } = styles;

ChartJS.register(ArcElement, Tooltip, Legend, ChartDataLabels);

export type Slice = {
  id: string | number;
  label: string;
  value: { count: number; percentage: number };
  color: string;
};

export type PieChartData = {
  label?: string;
  slices: Slice[];
};

type Props = {
  data: PieChartData;
  onSliceClick?: (id: string | number) => void;
};

export const PieChart = (props: Props) => {
  const { data, onSliceClick } = props;

  const chartRef = useRef<ChartJS<'pie'> | null>(null);

  const colors = data.slices.map((slice) => slice.color);

  const chartData = {
    labels: data.slices.map((child) => child.label),
    datasets: [
      {
        label: data.label || '',
        data: data.slices.map((child) => child.value.count),
        backgroundColor: colors,
        borderWidth: 1,
        borderColor: '#fff',
        hoverBorderWidth: 1,
        hoverBorderColor: '#fff',
      },
    ],
  };

  const handleClick = (event: React.MouseEvent<HTMLCanvasElement>) => {
    if (!chartRef.current || !onSliceClick) {
      return;
    }

    const points = chartRef.current.getElementsAtEventForMode(
      event.nativeEvent,
      'nearest',
      { intersect: true },
      true
    );

    if (!points.length) {
      return;
    }

    const clickedIndex = points[0].index;
    const clickedChild = data.slices[clickedIndex];

    if (clickedChild?.id != null) {
      onSliceClick(clickedChild.id);
    }
  };

  return (
    <div className={pie}>
      <Pie
        ref={chartRef}
        data={chartData}
        onClick={handleClick}
        options={{
          maintainAspectRatio: false,
          responsive: true,
          plugins: {
            tooltip: {
              enabled: true,
              mode: 'nearest',
              intersect: true,
              position: 'nearest',
              displayColors: false,
              backgroundColor: 'rgba(0,0,0,0.75)',
              padding: 8,
              callbacks: {
                title: (items) => {
                  const idx = items?.[0]?.dataIndex;
                  return typeof idx === 'number' ? data.slices[idx].label : '';
                },
                label: (ctx) => {
                  const idx = ctx.dataIndex;
                  const slice = data.slices[idx];
                  if (!slice) return '';
                  return `${slice.value.count} — ${slice.value.percentage}%`;
                },
              },
            },
            legend: { display: false },
            datalabels: {
              display: false,
            },
          },
          cutout: '35%',
          onHover: (event, activeElements) => {
            const canvas = event.native?.target as HTMLCanvasElement;
            canvas.style.cursor =
              activeElements.length > 0 ? 'pointer' : 'default';
          },
        }}
      />
    </div>
  );
};
