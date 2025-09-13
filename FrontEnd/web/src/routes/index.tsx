import { PieChart } from '../components/chart/PieChart';
import styles from './styles/index.module.css';

const { Div } = styles;

export const Index = () => {
  return (
    <div className={Div}>
      <PieChart isFirstLevel={true} />
    </div>
  );
};
