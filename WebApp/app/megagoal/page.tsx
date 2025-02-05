import { Suspense } from 'react';
import TeamsGrid from '@/app/ui/dashboard/teamsGrid';

export default async function Page() {
  return (
    <>
      <Suspense fallback={<p>Loading teams...</p>}>
        <TeamsGrid />
      </Suspense>
    </>
  );
}