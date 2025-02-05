export interface Team {
    team: {
        id: number;
        name: string;
        code: string;
        country: string;
        founded: number;
        national: boolean;
        logo: string;
    };
    venue: {
        id: number;
        name: string;
        address: string;
        city: string;
        capacity: number;
        surface: string;
        image: string;
    };
    seasons: {
        league: string;
        season: string;
    };
    previous: string[];
}