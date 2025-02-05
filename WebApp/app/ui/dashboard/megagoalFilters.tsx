'use client';
import { useState, useEffect, Key, Dispatch, SetStateAction } from 'react';
import { useSearchParams, usePathname, useRouter } from 'next/navigation';
import { Autocomplete, AutocompleteItem } from "@nextui-org/react";

import { useTopLeagues } from '@/app/lib/data';
import { League } from '@/app/modals/league';

interface MegagoalFiltersProps {
  order: number;
  setOrder: Dispatch<SetStateAction<number>>;
}

export default function MegagoalFilters({ order, setOrder }: MegagoalFiltersProps) {
  const searchParams = useSearchParams();
  const [countryFilter, setCountryFilter] = useState<string | null>(searchParams?.get('country') || null);
  const [leagueFilter, setLeagueFilter] = useState<string | null>(searchParams?.get('league_id') || null);
  const [seasonFilter, setSeasonFilter] = useState<string | null>(searchParams?.get('season') || null);
  const pathname = usePathname();
  const { replace } = useRouter();

  const { data: topLeagues } = useTopLeagues();

  useEffect(() => {
    const params = new URLSearchParams(searchParams?.toString());
      if (leagueFilter) {
        params.set('league_id', leagueFilter.toString());
      } else {
        params.delete('league_id');
      }
      if (countryFilter) {
        params.set('country', countryFilter.toString());
      } else {
        params.delete('country');
      }
      if (seasonFilter) {
        params.set('season', seasonFilter.toString());
      } else {
        params.delete('season');
      }
      replace(`${pathname}?${params.toString()}`);
  }, [countryFilter, leagueFilter, seasonFilter, pathname, replace, searchParams]);

  const countries = [
    { name: 'Spain' },
    { name: 'England' },
    { name: 'Italy' },
    { name: 'Germany' },
    { name: 'France' }
  ];

  const seasons = [
    { name: '2024/2025', value: '2024' },
    { name: '2023/2024', value: '2023' },
    { name: '2022/2023', value: '2022' },
    { name: '2021/2022', value: '2021' },
    { name: '2020/2021', value: '2020' },
    { name: '2019/2020', value: '2019' },
    { name: '2018/2019', value: '2018' },
    { name: '2017/2018', value: '2017' },
    { name: '2016/2017', value: '2016' },
    { name: '2015/2016', value: '2015' },
    { name: '2014/2015', value: '2014' }
  ];

  const orderList = [
    { name: '1 Up', value: 1 },
    { name: '1 Down', value: 2 },
    { name: 'A Up', value: 3 },
    { name: 'A Down', value: 4 },
  ];

  const onSelectionChangeCountryFilter = (key: Key | null) => {
    if (key !== null) {
      setCountryFilter(key.toString());
    } else {
      setCountryFilter(null);
    }
  };

  const onSelectionChangeLeagueFilter = (key: Key | null) => {
    if (key !== null) {
      setLeagueFilter(key.toString());
    } else {
      setLeagueFilter(null);
    }
  };
  
  const onSelectionChangeSeasonFilter = (key: Key | null) => {
    if (key !== null) {
      setSeasonFilter(key.toString());
    } else {
      setSeasonFilter(null);
    }
  };
  
  const onSelectionChangeOrderList = (key: Key | null) => {
      if (key !== null) {
          setOrder(Number(key));
      }
  };

  return (
    <>
      <div className='mb-5 gap-2 grid grid-cols-5'>

        <Autocomplete
          className="max-w-xs"
          defaultItems={countries}
          defaultSelectedKey={countryFilter || undefined}
          label="Search by country"
          onSelectionChange={onSelectionChangeCountryFilter}
        >
          {(country) => <AutocompleteItem key={country.name}>{country.name}</AutocompleteItem>}
        </Autocomplete>

        <Autocomplete
          className="max-w-xs"
          defaultItems={topLeagues}
          defaultSelectedKey={leagueFilter || undefined}
          label="Search by league"
          onSelectionChange={onSelectionChangeLeagueFilter}
        >
          {(topLeague: League) => <AutocompleteItem key={topLeague.league.id}>{topLeague.league.name}</AutocompleteItem>}
        </Autocomplete>

        <Autocomplete
          className="max-w-xs"
          defaultItems={seasons}
          defaultSelectedKey={seasonFilter || undefined}
          label="Search by season"
          onSelectionChange={onSelectionChangeSeasonFilter}
        >
          {(season) => <AutocompleteItem key={season.value}>{season.name}</AutocompleteItem>}
        </Autocomplete>

        <Autocomplete
          className="max-w-xs"
          defaultItems={orderList}
          defaultSelectedKey={order?.toString() || undefined}
          label="Order by"
          onSelectionChange={onSelectionChangeOrderList}
        >
          {(order) => <AutocompleteItem key={order.value}>{order.name}</AutocompleteItem>}
        </Autocomplete>

      </div>
    </>
  );
};

