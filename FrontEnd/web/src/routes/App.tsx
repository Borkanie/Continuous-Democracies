import styles from './styles/App.module.css';
import { Header } from '../components/header/Header';
import { RoundsList } from '../components/rounds-list/RoundsList';
import { RoundSection } from '../components/round-section/RoundSection';

const { wrapper } = styles;

export const App = () => {
  return (
    <>
      <Header />
      <div className={wrapper}>
        <RoundsList />
        <RoundSection />
      </div>
    </>
  );
};
