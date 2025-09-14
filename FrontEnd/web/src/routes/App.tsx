import { Outlet } from '@tanstack/react-router';
import styles from './styles/App.module.css';
import { Header } from '../components/header/Header';
import { RoundsList } from '../components/rounds-list/RoundsList';

const { wrapper } = styles;

export const App = () => {
  return (
    <>
      <Header />
      <div className={wrapper}>
        <RoundsList />
        <Outlet />
      </div>
    </>
  );
};
