'use client';

import MegagoalFilters from './megagoalFilters';
import MegagoalItems from './megagoalItems';
import { Suspense, useState } from 'react';
import {Divider} from "@nextui-org/react";

import { TeamsSkeleton } from '@/app/ui/skeletons';

export default function TeamsGrid() {
  const [order, setOrder] = useState(1);

  return (
    <>
      <MegagoalFilters order={order} setOrder={setOrder} />

      <Divider />

      <Suspense fallback={<TeamsSkeleton/>}>
        <MegagoalItems order={order} />
      </Suspense>
    </>
  );
};

