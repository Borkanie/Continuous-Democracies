import { useRef } from 'react';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Pie } from 'react-chartjs-2';
import ChartDataLabels from 'chartjs-plugin-datalabels';
import styles from './PieChart.module.css';
import { useNavigate } from '@tanstack/react-router';

const { pie } = styles;

ChartJS.register(ArcElement, Tooltip, Legend, ChartDataLabels);

type ChildEntry = {
  id: string;
  label: string;
  value: number;
};

type DrilldownEntry = {
  id: string;
  label: string;
  children: ChildEntry[];
};

const topLevel: DrilldownEntry = {
  id: '0',
  label: 'Vot',
  children: [
    { id: '1', label: 'Da', value: 8 },
    { id: '2', label: 'Nu', value: 6 },
    { id: '3', label: 'Absent', value: 8 },
    { id: '4', label: 'Abtinere', value: 6 },
  ],
};

const drilldownData: DrilldownEntry[] = [
  {
    id: '1',
    label: 'Da',
    children: [
      { id: '101', label: 'PSD', value: 3 },
      { id: '102', label: 'PNL', value: 2 },
      { id: '103', label: 'USR', value: 3 },
    ],
  },
  {
    id: '2',
    label: 'Nu',
    children: [
      { id: '101', label: 'PSD', value: 3 },
      { id: '102', label: 'PNL', value: 2 },
      { id: '103', label: 'USR', value: 3 },
    ],
  },
  {
    id: '4',
    label: 'Abtinere',
    children: [
      { id: '101', label: 'PSD', value: 3 },
      { id: '102', label: 'PNL', value: 2 },
      { id: '103', label: 'USR', value: 3 },
    ],
  },
  {
    id: '3',
    label: 'Absent',
    children: [
      { id: '101', label: 'PSD', value: 3 },
      { id: '102', label: 'PNL', value: 2 },
      { id: '103', label: 'USR', value: 3 },
    ],
  },
];

const COLORS = ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF'];

type Props = {
  isFirstLevel?: boolean;
  drilldownId?: string;
};

export const PieChart = (props: Props) => {
  const { isFirstLevel = false, drilldownId } = props;

  const chartRef = useRef<ChartJS<'pie'> | null>(null);
  const navigate = useNavigate();

  const current: DrilldownEntry = isFirstLevel
    ? topLevel
    : drilldownData.find((entry) => entry.id === drilldownId) ?? {
        id: '-1',
        label: '',
        children: [],
      };

  const chartData = {
    labels: current.children.map((child) => child.label),
    datasets: [
      {
        label: current.label,
        data: current.children.map((child) => child.value),
        backgroundColor: COLORS,
        borderWidth: 2,
        borderColor: '#fff',
        hoverBorderWidth: 2,
        hoverBorderColor: '#fff',
      },
    ],
  };

  const handleClick = (event: React.MouseEvent<HTMLCanvasElement>) => {
    if (!isFirstLevel) {
      return;
    }

    const chart = chartRef.current;
    if (!chart) {
      return;
    }

    const points = chart.getElementsAtEventForMode(
      event.nativeEvent,
      'nearest',
      { intersect: true },
      true
    );

    if (!points.length) {
      return;
    }

    const clickedIndex = points[0].index;
    const clickedChild = current.children[clickedIndex];

    if (clickedChild?.id != null) {
      navigate({ to: `section/${clickedChild.id}` });
    }
  };

  const handleBack = () => navigate({ to: '/' });

  return (
    <div>
      <h2>{isFirstLevel ? 'Vot' : `Subcategories of: ${current.label}`}</h2>
      {!isFirstLevel && <button onClick={handleBack}>⬅ Back</button>}
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
                enabled: false, // Disable default tooltip for sections
              },
              legend: {
                position: 'bottom',
                onClick: () => {}, // Disable legend click
              },
              datalabels: {
                color: '#fff',
                font: {
                  weight: 'bold',
                  size: 14,
                },
                // Show the label on the pie section instead of the tooltip
                formatter: (value, context) => {
                  const label = context.chart.data.labels?.[context.dataIndex];
                  return `${label}: ${value}`;
                },
              },
            },
            cutout: '35%',
            onHover: (event, activeElements) => {
              // Add cursor pointer on section hover.
              // event.native.target is the canvas element.
              const canvas = event.native?.target as HTMLCanvasElement;
              canvas.style.cursor =
                activeElements.length > 0 ? 'pointer' : 'default';
            },
          }}
        />
      </div>
    </div>
  );
};
