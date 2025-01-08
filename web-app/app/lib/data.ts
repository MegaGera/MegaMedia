import { useSuspenseQuery } from '@tanstack/react-query';

export async function fetchTopLeagues() {
  // await new Promise(resolve => setTimeout(resolve, 2000));
  const response = await fetch('http://localhost:3150/league/top/');
  if (!response.ok) {
    throw new Error('Failed to fetch top leagues.');
  }
  return response.json();
}

export async function fetchTeams(countryFilter: string | null, leagueFilter: string | null, seasonFilter: string | null) {
  // await new Promise(resolve => setTimeout(resolve, 2000));
  const params = new URLSearchParams();
  if (leagueFilter) {
    params.set('league_id', leagueFilter.toString());
  }
  if (countryFilter) {
    params.set('country', countryFilter.toString());
  }
  if (seasonFilter) {
    params.set('season', seasonFilter.toString());
  }

  const response = await fetch(`http://localhost:3150/team/?${params.toString()}`);
  if (!response.ok) {
    throw new Error('Failed to fetch teams.');
  }
  return response.json();
}

export function useTopLeagues() {
  const result = useSuspenseQuery({
    queryKey: ['topLeagues'],
    queryFn: fetchTopLeagues,
    suspense: true,
  });
  return result;
}

export function useTeams(countryFilter: string | null, leagueFilter: string | null, seasonFilter: string | null) {
  const result = useSuspenseQuery({
    queryKey: ['teams', countryFilter, leagueFilter, seasonFilter],
    queryFn: () => fetchTeams(countryFilter, leagueFilter, seasonFilter),
    suspense: true,
  });
  return result;
}

