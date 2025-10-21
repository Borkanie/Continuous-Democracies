import { useParams } from '@tanstack/react-router';

export const Section = () => {
  const { sectionId } = useParams({ strict: false });

  return <div>Section {sectionId}</div>;
};
