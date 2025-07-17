import { Outlet } from '@tanstack/react-router';
import styles from './styles/root.module.css';

const { wrapper } = styles;

export const Root = () => {
  return (
    <div className={wrapper}>
      <Outlet />
    </div>
  );
};
