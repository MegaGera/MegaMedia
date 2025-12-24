'use client';

import MegagoalFilters from './megagoalFilters';
import MegagoalItems from './megagoalItems';
import MegagoalLeaguesItems from './megagoalLeaguesItems';
import { Suspense, useState } from 'react';
import {Divider, Tabs, Tab} from "@nextui-org/react";

import { TeamsSkeleton } from '@/app/ui/skeletons';

export default function TeamsGrid() {
  const [order, setOrder] = useState(1);
  const [selectedView, setSelectedView] = useState('teams');

  return (
    <>
      <Tabs 
        selectedKey={selectedView} 
        onSelectionChange={(key) => setSelectedView(key as string)}
        className="mb-5"
      >
        <Tab key="teams" title="Teams" />
        <Tab key="leagues" title="Leagues" />
      </Tabs>

      {selectedView === 'teams' && (
        <>
          <MegagoalFilters order={order} setOrder={setOrder} />
          <Divider />
          <Suspense fallback={<TeamsSkeleton/>}>
            <MegagoalItems order={order} />
          </Suspense>
        </>
      )}

      {selectedView === 'leagues' && (
        <>
          <Divider />
          <Suspense fallback={<TeamsSkeleton/>}>
            <MegagoalLeaguesItems order={order} />
          </Suspense>
        </>
      )}
    </>
  );
};

