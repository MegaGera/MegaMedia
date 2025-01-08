'use client';

import MegagoalFilters from './megagoalFilters';
import MegagoalItems from './megagoalItems';
import { Suspense } from 'react';
import {Divider} from "@nextui-org/react";

import { TeamsSkeleton } from '@/app/ui/skeletons';

export default function TeamsGrid() {
  return (
    <>
      <MegagoalFilters />

      <Divider />

      <Suspense fallback={<TeamsSkeleton/>}>
        <MegagoalItems />
      </Suspense>
    </>
  );
};

