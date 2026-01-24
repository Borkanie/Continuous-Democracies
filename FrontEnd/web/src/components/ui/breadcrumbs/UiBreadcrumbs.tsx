import Breadcrumbs from '@mui/material/Breadcrumbs';
import { Link } from '@tanstack/react-router';
import styles from './UiBreadcrumbs.module.css';
import { UiText } from '../text/UiText';
import type { Breadcrumb } from '../../../utils/hooks/useBreadcrumbs';

const { link, current } = styles;

type Props = {
  breadcrumbs: Breadcrumb[];
};

export const UiBreadcrumbs = ({ breadcrumbs }: Props) => {
  return (
    <Breadcrumbs
      aria-label='breadcrumb'
      sx={{
        '& .MuiBreadcrumbs-separator': {
          color: 'var(--theme-font-color-secondary)',
          fontSize: 'var(--theme-font-size-m)',
        },
      }}
    >
      {breadcrumbs.map((crumb, index) => {
        const isLast = index === breadcrumbs.length - 1;

        if (isLast) {
          return (
            <UiText
              key={index}
              text={crumb.label}
              title={crumb.label}
              size={'medium'}
              className={current}
            />
          );
        }

        return (
          <Link key={index} to={crumb.to} className={link}>
            {crumb.label}
          </Link>
        );
      })}
    </Breadcrumbs>
  );
};
