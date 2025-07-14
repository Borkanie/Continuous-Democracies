import { useRef, useState } from 'react';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Pie } from 'react-chartjs-2';
import ChartDataLabels from 'chartjs-plugin-datalabels';
import styles from './PieChart.module.css';

const { pie } = styles;

ChartJS.register(ArcElement, Tooltip, Legend, ChartDataLabels);

type ChildEntry = {
  label: string;
  value: number;
};

type DrilldownEntry = {
  label: string;
  children: ChildEntry[];
};

const topLevel: DrilldownEntry = {
  label: 'Vot',
  children: [
    { label: 'Da', value: 8 },
    { label: 'Nu', value: 6 },
    { label: 'Absent', value: 8 },
    { label: 'Abtinere', value: 6 },
  ],
};

const drilldownData: DrilldownEntry[] = [
  {
    label: 'Da',
    children: [
      { label: 'PSD', value: 3 },
      { label: 'PNL', value: 2 },
      { label: 'USR', value: 3 },
    ],
  },
  {
    label: 'Nu',
    children: [
      { label: 'PSD', value: 3 },
      { label: 'PNL', value: 2 },
      { label: 'USR', value: 3 },
    ],
  },
  {
    label: 'Abtinere',
    children: [
      { label: 'PSD', value: 3 },
      { label: 'PNL', value: 2 },
      { label: 'USR', value: 3 },
    ],
  },
  {
    label: 'Absent',
    children: [
      { label: 'PSD', value: 3 },
      { label: 'PNL', value: 2 },
      { label: 'USR', value: 3 },
    ],
  },
];

const COLORS = ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF'];

export const PieChart = () => {
  const chartRef = useRef<ChartJS<'pie'> | null>(null);
  const [isFirstLevel, setIsFirstLevel] = useState(true);
  const [drilldownLabel, setDrilldownLabel] = useState<string | null>(null);

  const current: DrilldownEntry = isFirstLevel
    ? topLevel
    : drilldownData.find((entry) => entry.label === drilldownLabel) ?? {
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
        borderWidth: 1,
      },
    ],
  };

  const handleClick = (event: React.MouseEvent<HTMLCanvasElement>) => {
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
    const clickedLabel = chartData.labels?.[clickedIndex];

    if (isFirstLevel) {
      const match = drilldownData.find((entry) => entry.label === clickedLabel);
      if (match) {
        setDrilldownLabel(clickedLabel);
        setIsFirstLevel(false);
      }
    }
  };

  const handleBack = () => {
    setIsFirstLevel(true);
    setDrilldownLabel(null);
  };

  return (
    <div>
      <h2>{isFirstLevel ? 'Vot' : `Subcategories of ${drilldownLabel}`}</h2>
      {!isFirstLevel && <button onClick={handleBack}>⬅ Back</button>}
      <div className={pie}>
        <Pie
          ref={chartRef}
          data={chartData}
          onClick={handleClick}
          options={{
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

              if (activeElements.length > 0) {
                canvas.style.cursor = 'pointer';
              } else {
                canvas.style.cursor = 'default';
              }
            },
          }}
        />
      </div>
    </div>
  );
};
