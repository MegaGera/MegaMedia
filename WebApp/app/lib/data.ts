import { useSuspenseQuery } from '@tanstack/react-query';
import { CONFIG } from '@/app/constants';

// Define a helper function to add default headers and credentials to fetch calls
async function fetchWithHeaders(url: string, options: RequestInit = {}) {
  const defaultHeaders = {
    'Content-Type': 'application/json',
  };

  return fetch(url, {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers,
    },
    credentials: CONFIG.nodeEnv === "production" ? 'include' : undefined, // This ensures cookies are sent with the request
  });
}

export async function fetchTopLeagues() {
  console.log(`${CONFIG.megagoalServerApiUrl}`);
  const response = await fetchWithHeaders(`${CONFIG.megagoalServerApiUrl}/league/top/`);
  if (!response.ok) {
    throw new Error('Failed to fetch top leagues.');
  }
  return response.json();
}

export async function fetchTeams(countryFilter: string | null, leagueFilter: string | null, seasonFilter: string | null) {
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

  const response = await fetchWithHeaders(`${CONFIG.megagoalServerApiUrl}/team/?${params.toString()}`);
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

export function uploadTeamImage(teamID: string, image: File) {
  const formData = new FormData();
  formData.append('image', image);
  return fetchWithHeaders(`${CONFIG.megamediaServerApiUrl}/megagoal/team/${teamID}/image/`, {
    method: 'POST',
    body: formData,
    headers: {}, // Do not include 'Content-Type', as it is automatically set for FormData
  });
}

export function squaredTeamImage(teamID: string) {
  return fetchWithHeaders(`${CONFIG.megamediaServerApiUrl}/megagoal/team/${teamID}/squared/`, {
    method: 'POST',
  });
}

export function deleteTeamImage(teamID: string) {
  return fetchWithHeaders(`${CONFIG.megamediaServerApiUrl}/megagoal/team/${teamID}/delete/`, {
    method: 'POST',
  });
}
