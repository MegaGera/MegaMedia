'use client'

import { Suspense } from 'react';
import MegageraGrid from '../ui/dashboard/megageraGrid';

export default function Page() {

  return (
    <>
      <Suspense fallback={<p>Loading teams...</p>}>
        <MegageraGrid />
      </Suspense>
    </>
  );
}