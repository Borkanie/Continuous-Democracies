import { Outlet } from '@tanstack/react-router';
import styles from './OutletSection.module.css';

const { Div } = styles;

export const OutletSection = () => {
  return (
    <div className={Div}>
      <Outlet />
    </div>
  );
};
