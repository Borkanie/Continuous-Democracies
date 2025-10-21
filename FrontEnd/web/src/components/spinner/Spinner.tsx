import styles from './Spinner.module.css';

const { wrapper, spinner } = styles;

export const Spinner = () => (
  <div className={wrapper}>
    <div className={spinner}></div>
  </div>
);
