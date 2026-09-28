import ToggleButtonGroup from '@mui/material/ToggleButtonGroup';
import ToggleButton from '@mui/material/ToggleButton';
import { motion } from 'motion/react';
import { useChamber } from '../../../utils/context/ChamberContext';
import type { Chamber } from '../../../utils/types';

export const ChamberToggle = () => {
  const { chamber, toggleChamber } = useChamber();
  const isSenate = chamber === 'senate';

  const handleChange = (
    _: React.MouseEvent<HTMLElement>,
    next: Chamber | null,
  ) => {
    if (next && next !== chamber) {
      toggleChamber();
    }
  };

  return (
    <div
      style={{
        position: 'relative',
        display: 'inline-flex',
        padding: 3,
        border: 'var(--theme-border-primary)',
        backgroundColor: 'var(--background-color-tertiary)',
        borderRadius: 999,
        overflow: 'hidden',
      }}
    >
      <motion.span
        layout
        transition={{ duration: 0.35, ease: 'easeInOut' }}
        style={{
          position: 'absolute',
          top: 3,
          bottom: 3,
          left: isSenate ? '50%' : 3,
          width: 'calc(50% - 3px)',
          borderRadius: 999,
          backgroundColor: 'var(--background-color-secondary)',
        }}
      />
      <ToggleButtonGroup
        value={chamber}
        exclusive
        onChange={handleChange}
        aria-label={'Comuta intre Camera Deputatilor si Senat'}
      >
        <ToggleButton value={'parliament' as Chamber}>
          Camera Deputatilor
        </ToggleButton>
        <ToggleButton value={'senate' as Chamber}>Senat</ToggleButton>
      </ToggleButtonGroup>
    </div>
  );
};
