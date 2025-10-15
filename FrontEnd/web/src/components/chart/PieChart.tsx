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
};

export type PieChartData = {
  label?: string;
  slices: Slice[];
};

type Props = {
  data: PieChartData;
  onSliceClick?: (id: string | number) => void;
};
const COLORS = ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF'];

export const PieChart = (props: Props) => {
  const { data, onSliceClick } = props;

  const chartRef = useRef<ChartJS<'pie'> | null>(null);

  const chartData = {
    labels: data.slices.map((child) => child.label),
    datasets: [
      {
        label: data.label || '',
        data: data.slices.map((child) => child.value.count), // 👈 use count
        backgroundColor: COLORS,
        borderWidth: 2,
        borderColor: '#fff',
        hoverBorderWidth: 2,
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
            tooltip: { enabled: false },
            legend: { display: false },
            datalabels: {
              color: '#fff',
              font: {
                family: 'Montserrat, sans-serif',
                weight: 'bold',
                size: 14,
              },
              formatter: (value, context) => {
                const slice = data.slices[context.dataIndex];
                if (slice.value.count <= 0) {
                  return null;
                }
                const label = slice.label;
                return `${label}: ${slice.value.count} (${slice.value.percentage}%)`;
              },
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
