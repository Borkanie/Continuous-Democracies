import { faSearch } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import styles from './Search.module.css';

const { Div, iconContainer, icon } = styles;

export const Search = () => {
  return (
    <div className={Div}>
      <div className={iconContainer}>
        <FontAwesomeIcon icon={faSearch} className={icon} />
      </div>
      <input placeholder={'Cauta...'} />
    </div>
  );
};
