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
        animate={{ x: isSenate ? 'calc(100% + 5px)' : 0 }}
        transition={{ duration: 0.35, ease: 'easeInOut' }}
        style={{
          position: 'absolute',
          top: 3,
          bottom: 3,
          left: 3,
          width: 'calc(50% - 5.5px)',
          borderRadius: 999,
          background:
            'color-mix(in srgb, var(--theme-color-primary) 22%, var(--background-color-primary) 78%)',
          border: 'var(--theme-border-secondary)',
          boxShadow: 'var(--theme-box-shadow-primary)',
          backdropFilter: 'blur(6px)',
          WebkitBackdropFilter: 'blur(6px)',
        }}
      />
      <ToggleButtonGroup
        value={chamber}
        exclusive
        onChange={handleChange}
        aria-label={'Comuta intre Camera Deputatilor si Senat'}
        sx={{ gap: '5px' }}
      >
        <ToggleButton
          value={'parliament' as Chamber}
          sx={{ flex: 1, paddingInline: '19px' }}
        >
          Camera Deputatilor
        </ToggleButton>
        <ToggleButton
          value={'senate' as Chamber}
          sx={{ flex: 1, paddingInline: '19px' }}
        >
          Senat
        </ToggleButton>
      </ToggleButtonGroup>
    </div>
  );
};
