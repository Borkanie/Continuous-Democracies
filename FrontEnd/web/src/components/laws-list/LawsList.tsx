import { LawCard } from '../law-card/LawCard';
import { Search } from '../search/Search';
import styles from './LawsList.module.css';

const { Div, header } = styles;

export const LawsList = () => {
  return (
    <div className={Div}>
      <div className={header}>
        <h2>Lista legi</h2>
        <p>244</p>
      </div>
      <Search />
      <LawCard />
      <LawCard />
      <LawCard />
    </div>
  );
};
