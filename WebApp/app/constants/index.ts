export const CONFIG = {
    megagoalServerApiUrl: process.env.NEXT_PUBLIC_MEGAGOAL_SERVER_API_URL || 'http://localhost:3150',
    megamediaServerApiUrl: process.env.NEXT_PUBLIC_MEGAMEDIA_SERVER_API_URL || 'http://localhost:8080',
    nodeEnv: process.env.NODE_ENV || 'development'
};
