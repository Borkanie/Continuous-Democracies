import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSun, faMoon } from '@fortawesome/free-solid-svg-icons';
import { useTheme } from '../../../utils/context/ThemeContext';
import styles from './ThemeToggle.module.css';

export const ThemeToggle = () => {
  const { theme, toggleTheme } = useTheme();
  const isDark = theme === 'dark';

  return (
    <button
      className={styles.button}
      onClick={toggleTheme}
      aria-label={'Schimba tema'}
      title={'Schimba tema'}
    >
      <span
        className={`${styles.iconWrapper} ${
          isDark ? styles.dark : styles.light
        }`}
      >
        <FontAwesomeIcon icon={faSun} />
        <FontAwesomeIcon icon={faMoon} />
      </span>
    </button>
  );
};
