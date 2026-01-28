import { faChartSimple } from '@fortawesome/free-solid-svg-icons';
import styles from './Header.module.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { UiText } from '../text/UiText';
import { ThemeToggle } from '../theme-toggle/ThemeToggle';

const { header, content, iconContainer, icon, details, flex } = styles;

export const Header = () => {
  const heading = 'Voturi parlamentare';
  const subHeading = 'Vizualizare in timp real a voturilor legislative';

  return (
    <header className={header}>
      <div className={content}>
        <div>
          <div className={flex}>
            <div className={iconContainer}>
              <FontAwesomeIcon icon={faChartSimple} className={icon} />
            </div>
            <div className={details}>
              <h1>{heading}</h1>
              <UiText text={subHeading} title={subHeading} />
            </div>
          </div>
        </div>
        <ThemeToggle />
      </div>
    </header>
  );
};
