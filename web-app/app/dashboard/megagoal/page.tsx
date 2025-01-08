'use client';

import { useState, useEffect, Key } from 'react';
import { useSearchParams, usePathname, useRouter } from 'next/navigation';
import Link from 'next/link';

import {Card, CardBody, CardFooter, Image, Autocomplete, AutocompleteItem, Divider} from "@nextui-org/react";
import { use } from 'framer-motion/m';

export default function Page() {
  const searchParams = useSearchParams()
 
  const [topLeagues, setTopLeagues] = useState([]);
  const [teams, setTeams] = useState([]);
  const [countryFilter, setCountryFilter] = useState<string | null>(searchParams?.get('country') || null);
  const [leagueFilter, setLeagueFilter] = useState<number | string | null>(searchParams?.get('league_id') || null);
  const [seasonFilter, setSeasonFilter] = useState<string | null>(searchParams?.get('season') || null);
  const pathname = usePathname();
  const { replace } = useRouter();

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

  useEffect(() => {
    const fetchTopLeagues = async () => {
      const response = await fetch('http://localhost:3150/league/top/');
      const data = await response.json();
      setTopLeagues(data);
    };

    fetchTopLeagues();
  }, []);

  useEffect(() => {
    console.log("fetch teams")
    const fetchTeams = async () => {

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
      console.log("params", params.toString());

      replace(`${pathname}?${params.toString()}`);

      const response = await fetch(`http://localhost:3150/team/?${params.toString()}`);
      const data = await response.json();
      setTeams(data);
    };

    fetchTeams();
  }, [countryFilter, leagueFilter, seasonFilter]);

  const getSrcImageByTeamID = (teamID: string) => {
    return `https://megamedia.megagera.com/megagoal/teams/team_${teamID}.png`;
  }  

  const onSelectionChangeCountryFilter = (key: Key | null) => {
    if (key !== null) {
      setCountryFilter(key.toString());
    } else {
      setCountryFilter(null);
    }
  };

  const onSelectionChangeLeagueFilter = (key: Key | null) => {
    if (key !== null) {
      setLeagueFilter(Number(key));
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

  // const onInputChange = (value: any) => {
  //   console.log("input", value);
  // };

  return <>
    <div className='mb-5 gap-2 grid grid-cols-5'>

      <Autocomplete
        className="max-w-xs"
        defaultItems={countries}
        defaultSelectedKey={searchParams?.get('country') || undefined}
        label="Search by country"
        placeholder="Select a country"
        // onInputChange={onInputChange}
        onSelectionChange={onSelectionChangeCountryFilter}
      >
        {(country) => <AutocompleteItem key={country.name}>{country.name}</AutocompleteItem>}
      </Autocomplete>

      <Autocomplete
        className="max-w-xs"
        defaultItems={topLeagues}
        defaultSelectedKey={searchParams?.get('league_id') || undefined}
        label="Search by league"
        placeholder="Select a league"
        // onInputChange={onInputChange}
        onSelectionChange={onSelectionChangeLeagueFilter}
      >
        {(topLeague) => <AutocompleteItem key={topLeague['league']['id']}>{topLeague['league']['name']}</AutocompleteItem>}
      </Autocomplete>
      
      <Autocomplete
        className="max-w-xs"
        defaultItems={seasons}
        defaultSelectedKey={searchParams?.get('season') || undefined}
        label="Search by season"
        placeholder="Select a season"
        // onInputChange={onInputChange}
        onSelectionChange={onSelectionChangeSeasonFilter}
      >
        {(season) => <AutocompleteItem key={season.value}>{season.name}</AutocompleteItem>}
      </Autocomplete>

    </div>
    <Divider />
    <div className="mt-5 gap-2 grid grid-cols-3 sm:grid-cols-10">

      {teams.slice(0, 100).map((team, index) => (

        /* eslint-disable no-console */
        // <Link
        //     key={item.title}
        //     href={item.href}
        //     className='w-full'
        //   >
        <Card className='w-[90%] m-auto h-full' key={index} isPressable shadow="sm" onPress={() => console.log("item pressed")}>
          <CardBody className="overflow-visible p-0 m-auto">
            <Image
              alt={team['team']['name']}
              className="w-full object-cover h-full p-2"
              radius="lg"
              shadow="sm"
              src={getSrcImageByTeamID(team['team']['id'])}
              width="100%" />
          </CardBody>
          <CardFooter className="text-small grid">
            <b>{team['team']['name']}</b>
            <p className="text-default-500">{team['team']['id']}</p>
          </CardFooter>
        </Card>
        // </Link>
      ))}
    </div>
  </>
}