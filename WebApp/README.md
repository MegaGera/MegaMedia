# MegaMedia - WebApp

WebApp microservice built with [Next.js](https://nextjs.org/) (React) for *MegaMedia*.

## Table of Contents

- [Service Description](#service-description)
- [Project Composition](#project-composition)
- [Dependencies](#dependencies)
- [Environment Variables](#environment-variables)
- [Development Server](#development-server)
- [Build & Deploy](#build--deploy)

## Service Description

Web Application microservice for *MegaMedia* to handle the images of the rest of *Mega* services.

There are common features for all the services like update images or delete the previous ones, but for each microservice there could be individual features like Square the *MegaGoal* images. Also the filters to search through the static files of each microservice are different, having to connect with the Server of these Services (for now only with *MegaGoal* but it will be expanded).

UI responsive design. The design is adapted to be opened with big screens and mobile screens.

![megamedia gif](static/screenshots/megamedia.gif)

The application is already deployed and available at: [https://megamedia.megagera.com](https://megamedia.megagera.com)

The deployment in production is done with [Docker](https://www.docker.com/).

## Project composition

This project is developed with TypeScript files. It's composed with the next modules: 

- **Modals**:
  - Validate the data models of the data retrieved from the API. Located in [`app/modals`](app/modals).

- **Data**:
  - Service file that triggers HTTP calls to *MegaMedia Server* and *MegaGoal Server*. Located in [`app/lib/data.ts`](app/lib/data.ts).

- **UI Components**:
  - UI Components for the different pages, cards, grid, items, filters...

- **MegaGera and MegaGoal Pages**:
  - Main pages for visualize static files of both services.

The source code is under the [`app`](app) folder.

## Dependencies

The dependencies and the scripts are defined in the file [`package.json`](package.json) and managed with [pnpm](https://www.npmjs.com/).

To install the dependencies run the command: `pnpm install`.

## Environment Variables

In: `.env.development` || `.env.production`

```javascript
NEXT_PUBLIC_MEGAGOAL_SERVER_API_URL=string
NEXT_PUBLIC_MEGAMEDIA_SERVER_API_URL=string
```

## Development server

Run `pnpm start:dev` for a dev server. Navigate to `http://localhost:4200/`. The application will automatically reload if detects any change of the source files.

## Build & Deploy

[`Dockerfile`](Dockerfile) file builds the app for production and generates de Docker container.

[`docker-compose.yml`](docker-compose.yml) file manages the image and handle it easily within the *Mega* network.