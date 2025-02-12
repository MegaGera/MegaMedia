import { useSuspenseQuery } from '@tanstack/react-query';
import { CONFIG } from '@/app/constants';

// Define a helper function to add default headers and credentials to fetch calls
async function fetchWithHeaders(url: string, image: boolean, options: RequestInit = {}) {
  // const contentType = !image ? 'application/json' : 'multipart/form-data';
  const defaultHeaders: HeadersInit = !image ? {
    'Content-Type': 'application/json',
  } : {}

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
  const response = await fetchWithHeaders(`${CONFIG.megagoalServerApiUrl}/league/top/`, false);
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

  const response = await fetchWithHeaders(`${CONFIG.megagoalServerApiUrl}/team?${params.toString()}`, false);
  if (!response.ok) {
    throw new Error('Failed to fetch teams.');
  }
  return response.json();
}

export function useTopLeagues() {
  const result = useSuspenseQuery({
    queryKey: ['topLeagues'],
    queryFn: fetchTopLeagues,
    // suspense: true,
  });
  return result;
}

export function useTeams(countryFilter: string | null, leagueFilter: string | null, seasonFilter: string | null) {
  const result = useSuspenseQuery({
    queryKey: ['teams', countryFilter, leagueFilter, seasonFilter],
    queryFn: () => fetchTeams(countryFilter, leagueFilter, seasonFilter),
    // suspense: true,
  });
  return result;
}

export function uploadTeamImage(teamID: string, image: FormData) {
  return fetchWithHeaders(`${CONFIG.megamediaServerApiUrl}/api/megagoal/team/${teamID}/image/`, true, {
    method: 'POST',
    body: image,
    headers: {}, // Do not include 'Content-Type', as it is automatically set for FormData
  });
}

export function squaredTeamImage(teamID: string) {
  return fetchWithHeaders(`${CONFIG.megamediaServerApiUrl}/api/megagoal/team/${teamID}/squared/`, false, {
    method: 'POST',
  });
}

export function deleteTeamImage(teamID: string) {
  return fetchWithHeaders(`${CONFIG.megamediaServerApiUrl}/api/megagoal/team/${teamID}/delete/`, false, {
    method: 'POST',
  });
}

export function fetchMegageraLogos() {
  return fetchWithHeaders(`${CONFIG.megamediaServerApiUrl}/api/megagera/`, false, {
    method: 'GET',
  });
}
export function uploadMegageraImage(id: string, image: FormData) {
  return fetchWithHeaders(`${CONFIG.megamediaServerApiUrl}/api/megagera/${id}/upload/`, true, {
    method: 'POST',
    body: image,
    headers: {}, // Do not include 'Content-Type', as it is automatically set for FormData
  });
}

export function deleteMegageraImage(id: string, name: string) {
  return fetchWithHeaders(`${CONFIG.megamediaServerApiUrl}/api/megagera/${id}/delete/${name}/`, false, {
    method: 'POST',
  });
}
