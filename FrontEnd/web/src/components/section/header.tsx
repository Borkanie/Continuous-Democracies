import styles from './header.module.css';
import { DateComponent } from '../date/DateComponent';
import { UiText } from '../ui/text/UiText';
import { UiBreadcrumbs } from '../ui/breadcrumbs/UiBreadcrumbs';
import { useBreadcrumbs } from '../../utils/hooks/useBreadcrumbs';
import {
  faBarsStaggered,
  faCircleCheck,
  faCircleXmark,
} from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { UiButton } from '../ui/button/UiButton';
import { useBreakpoints } from '../../utils/hooks/useBreakpoints';
import { useDrawer } from '../../utils/context/DrawerContext';
import { useTheme } from '../../utils/context/ThemeContext';
import { positionColor, majorityTypeLabel } from '../../utils/helper';
import type { MajorityType } from '../../utils/types';

const { header, extraDetails, titleRow, statusGroup, statusItem, statusItemValue, statusLabel } =
  styles;

type Props = {
  title: string;
  subtitle?: string;
  description?: string;
  status?: 'PASSED' | 'FAILED';
  majorityType?: MajorityType;
  extraDetails?: {
    voteDate?: Date;
  };
};

export const Header = (props: Props) => {
  const {
    title: headerTitle,
    description,
    status,
    majorityType,
    extraDetails: details,
  } = props;

  const breadcrumbs = useBreadcrumbs();
  const { isDesktop } = useBreakpoints();
  const { setIsDrawerOpen } = useDrawer();
  const { theme } = useTheme();

  return (
    <div className={header}>
      {breadcrumbs.length > 1 && (
        <UiBreadcrumbs
          breadcrumbs={breadcrumbs.length > 1 ? breadcrumbs : []}
        />
      )}
      <div className={titleRow}>
        <h3>{headerTitle}</h3>
        {status && (
          <div className={statusGroup}>
            {majorityType && (
              <div className={statusItem}>
                <UiText
                  className={statusLabel}
                  text={'Modalitate de votare'}
                  size={'small'}
                />
                <UiText text={majorityTypeLabel(majorityType)} size={'small'} />
              </div>
            )}
            <div className={statusItem}>
              <UiText className={statusLabel} text={'A trecut'} size={'small'} />
              <div className={statusItemValue}>
                <UiText
                  text={status === 'PASSED' ? 'Da' : 'Nu'}
                  size={'small'}
                />
                <FontAwesomeIcon
                  icon={status === 'PASSED' ? faCircleCheck : faCircleXmark}
                  title={status}
                  style={{
                    color: positionColor(status === 'PASSED' ? 'Yes' : 'No', theme),
                  }}
                />
              </div>
            </div>
          </div>
        )}
      </div>
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
