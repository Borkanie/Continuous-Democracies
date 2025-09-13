import { Outlet } from '@tanstack/react-router';
import styles from './styles/root.module.css';
import { Header } from '../components/header/Header';
import { LawsList } from '../components/laws-list/LawsList';

const { wrapper } = styles;

export const Root = () => {
  return (
    <>
      <Header />
      <div className={wrapper}>
        <LawsList />
        <Outlet />
      </div>
    </>
  );
};
