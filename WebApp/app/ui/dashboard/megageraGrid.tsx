'use client';

import MegageraItems from './megageraItems';
import { Suspense } from 'react';

import { TeamsSkeleton } from '@/app/ui/skeletons';

export default function MegageraGrid() {

  return (
    <>
      <Suspense fallback={<TeamsSkeleton/>}>
        <MegageraItems />
      </Suspense>
    </>
  );
};

