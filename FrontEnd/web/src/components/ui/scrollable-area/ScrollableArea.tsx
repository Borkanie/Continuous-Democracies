import classNames from 'classnames';
import styles from './ScrollableArea.module.css';

const { outer, inner } = styles;

export const ScrollableArea = ({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) => {
  return (
    <div className={classNames(className, 'ScrollableArea', outer)}>
      <div className={inner}>{children}</div>
    </div>
  );
};
