import { useParams } from '@tanstack/react-router';
import { PieChart } from '../components/chart/PieChart';

export const Section = () => {
  const { sectionId } = useParams({ strict: false });

  return <PieChart drilldownId={sectionId} />;
};
