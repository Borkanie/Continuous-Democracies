import { createContext, useContext, useEffect, useState } from 'react';
import type { Chamber } from '../types';

const ChamberContext = createContext<{
  chamber: Chamber;
  toggleChamber: () => void;
}>({
  chamber: 'parliament',
  toggleChamber: () => {},
});

export function ChamberProvider({ children }: { children: React.ReactNode }) {
  const [chamber, setChamber] = useState<Chamber>(() => {
    return (localStorage.getItem('chamber') as Chamber) || 'parliament';
  });

  useEffect(() => {
    localStorage.setItem('chamber', chamber);
  }, [chamber]);

  const toggleChamber = () =>
    setChamber((current) =>
      current === 'parliament' ? 'senate' : 'parliament',
    );

  return (
    <ChamberContext.Provider value={{ chamber, toggleChamber }}>
      {children}
    </ChamberContext.Provider>
  );
}

export const useChamber = () => useContext(ChamberContext);
