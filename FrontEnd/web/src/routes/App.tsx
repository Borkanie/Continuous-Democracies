import styles from './styles/App.module.css';
import { Header } from '../components/header/Header';
import { RoundsList } from '../components/rounds-list/RoundsList';
import { useBreakpoints } from '../utils/hooks/useBreakpoints';
import { useDrawer } from '../utils/context/DrawerContext';
import Drawer from '@mui/material/Drawer';
import { OutletSection } from '../components/outlet-section/OutletSection';

const { wrapper } = styles;

export const App = () => {
  const { isDesktop } = useBreakpoints();
  const { isDrawerOpen, setIsDrawerOpen } = useDrawer();

  return (
    <>
      <Header />
      <div className={wrapper}>
        {isDesktop && <RoundsList />}
        <OutletSection />
      </div>
      {!isDesktop && (
        <Drawer
          open={isDrawerOpen}
          onClose={() => setIsDrawerOpen(false)}
          transitionDuration={300}
          slotProps={{
            paper: {
              sx: {
                width: '100%',
              },
            },
          }}
        >
          <RoundsList inDrawer={true} />
        </Drawer>
      )}
    </>
  );
};
