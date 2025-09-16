import { Outlet } from '@tanstack/react-router';
import styles from './RoundSection.module.css';

const { Div } = styles;

export const RoundSection = () => {
  return (
    <div className={Div}>
      <Outlet />
    </div>
  );
};
