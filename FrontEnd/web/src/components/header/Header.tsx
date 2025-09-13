import { faChartSimple } from '@fortawesome/free-solid-svg-icons';
import styles from './Header.module.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

const { header, content, iconContainer, icon, details } = styles;

export const Header = () => {
  return (
    <header className={header}>
      <div className={content}>
        <div className={iconContainer}>
          <FontAwesomeIcon icon={faChartSimple} className={icon} />
        </div>
        <div className={details}>
          <h1>Voturi parlamentare</h1>
          <p>Vizualizare in timp real a voturilor legislative</p>
        </div>
      </div>
    </header>
  );
};
