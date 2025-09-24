'use client';

import MegageraItems from './megageraItems';
import MegageraTopBar from './megageraTopBar';
import { Suspense, useState } from 'react';
import { Divider } from "@nextui-org/react";

import { TeamsSkeleton } from '@/app/ui/skeletons';

export default function MegageraGrid() {
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleImageAdded = () => {
    setRefreshTrigger(prev => prev + 1);
  };

  return (
    <>
      <MegageraTopBar onImageAdded={handleImageAdded} />
      
      <Divider />
      
      <Suspense fallback={<TeamsSkeleton/>} key={refreshTrigger}>
        <MegageraItems />
      </Suspense>
    </>
  );
};

