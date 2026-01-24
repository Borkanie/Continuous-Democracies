import { faChartSimple } from '@fortawesome/free-solid-svg-icons';
import styles from './Header.module.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { UiText } from '../ui/text/UiText';

const { header, content, iconContainer, icon, details } = styles;

export const Header = () => {
  const heading = 'Voturi parlamentare';
  const subHeading = 'Vizualizare in timp real a voturilor legislative';

  return (
    <header className={header}>
      <div className={content}>
        <div className={iconContainer}>
          <FontAwesomeIcon icon={faChartSimple} className={icon} />
        </div>
        <div className={details}>
          <h1>{heading}</h1>
          <UiText text={subHeading} title={subHeading} />
        </div>
      </div>
    </header>
  );
};
