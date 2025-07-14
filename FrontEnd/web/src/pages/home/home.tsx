import { PieChart } from '../../components/chart/PieChart';
import styles from './home.module.css';

const { wrapper } = styles;

export const Home = () => {
  return (
    <div className={wrapper}>
      <h1>Welcome to the Home Page</h1>
      <PieChart />
    </div>
  );
};
