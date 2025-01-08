'use client';
import { useSearchParams } from 'next/navigation';
import { Card, CardBody, CardFooter, Image } from "@nextui-org/react";

import { useTeams } from '@/app/lib/data';

export default function MegagoalItems() {
  const searchParams = useSearchParams();
  const countryFilter = searchParams?.get('country') || null;
  const leagueFilter = searchParams?.get('league_id') || null;
  const seasonFilter = searchParams?.get('season') || null;

  const { data: teams } = useTeams(countryFilter, leagueFilter, seasonFilter);

  const getSrcImageByTeamID = (teamID: string) => {
    return `https://megamedia.megagera.com/megagoal/teams/team_${teamID}.png`;
  }  
  
  return (
  <>
    <div className="mt-5 gap-2 grid grid-cols-3 sm:grid-cols-10">
      {teams.slice(0, 20).map((team, index) => (
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
      ))}
    </div>
    </>
  );
};

