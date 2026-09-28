import { AnimatePresence, motion } from 'motion/react';
import IconButton from '@mui/material/IconButton';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSun, faMoon } from '@fortawesome/free-solid-svg-icons';
import { useTheme } from '../../../utils/context/ThemeContext';

export const ThemeToggle = () => {
  const { theme, toggleTheme } = useTheme();
  const isDark = theme === 'dark';

  return (
    <IconButton
      onClick={toggleTheme}
      aria-label={'Schimba tema'}
      title={'Schimba tema'}
      sx={{ width: 32, height: 32, overflow: 'hidden' }}
    >
      <AnimatePresence mode={'wait'} initial={false}>
        <motion.span
          key={isDark ? 'moon' : 'sun'}
          initial={{ y: isDark ? 20 : -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          exit={{ y: isDark ? -20 : 20, opacity: 0 }}
          transition={{ duration: 0.35, ease: 'easeInOut' }}
          style={{ display: 'grid', placeItems: 'center' }}
        >
          <FontAwesomeIcon icon={isDark ? faMoon : faSun} />
        </motion.span>
      </AnimatePresence>
    </IconButton>
  );
};
