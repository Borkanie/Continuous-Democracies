import { Outlet } from '@tanstack/react-router';
import styles from './styles/root.module.css';
import { Header } from '../components/header/Header';

const { wrapper } = styles;

export const Root = () => {
  return (
    <>
      <Header />
      <div className={wrapper}>
        <Outlet />
      </div>
    </>
  );
};
