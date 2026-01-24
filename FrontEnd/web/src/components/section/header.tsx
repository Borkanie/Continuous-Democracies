import styles from './header.module.css';
import { DateComponent } from '../date/DateComponent';
import { UiText } from '../ui/text/UiText';
import { UiBreadcrumbs } from '../ui/breadcrumbs/UiBreadcrumbs';
import { useBreadcrumbs } from '../../utils/hooks/useBreadcrumbs';
import { faBarsStaggered } from '@fortawesome/free-solid-svg-icons';
import { UiButton } from '../ui/button/UiButton';
import { useBreakpoints } from '../../utils/hooks/useBreakpoints';
import { useDrawer } from '../../utils/context/DrawerContext';

const { header, extraDetails } = styles;

type Props = {
  title: string;
  subtitle?: string;
  description?: string;
  extraDetails?: {
    voteDate?: Date;
  };
};

export const Header = (props: Props) => {
  const { title: headerTitle, description, extraDetails: details } = props;

  const breadcrumbs = useBreadcrumbs();
  const { isDesktop } = useBreakpoints();
  const { setIsDrawerOpen } = useDrawer();

  return (
    <div className={header}>
      {breadcrumbs.length > 1 && (
        <UiBreadcrumbs
          breadcrumbs={breadcrumbs.length > 1 ? breadcrumbs : []}
        />
      )}
      <h3>{headerTitle}</h3>
      {description && <UiText text={description} />}
      <div className={extraDetails}>
        {!isDesktop && (
          <UiButton
            icon={faBarsStaggered}
            title={'Deschide lista legi'}
            text={'Legi'}
            onClick={() => setIsDrawerOpen(true)}
          />
        )}
        {details?.voteDate && <DateComponent text={details.voteDate} />}
      </div>
    </div>
  );
};
