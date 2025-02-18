# MegaMedia - Server

API Server built with [Go](https://go.dev/) for *MegaMedia*.

## Table of Contents

- [Service Description](#service-description)
- [Project Composition](#project-composition)
- [Dependencies](#dependencies)
- [Environment Variables](#environment-variables)
- [Development Server](#development-server)
- [Build & Deploy](#build--deploy)

## Service Description

Server *serves static files* to other *Mega* services.

API Server that exposes endpoints to handle the images for *MegaGoal* and *Mega*. The idea is for each  *Mega* service there would be different types of features to handle these static files, but reuse the modules as much as possible.

Features for each microservice:

**MegaGoal**

- Update team logo.
- See previous team logo.
- Delete previous team logo.
- Squared team logo. Make an image same size of height and width, to make them look better in *MegaGoal WebApp*.

After update a logo of a team, the *MegaGoal Database* is updated through the *MegaGoal Server*. It updates the metadata of the teams to keep the names of the previous images.

**Mega (MegaGera)**

- Update image.
- See previous image.
- Delete previous image.

After update an image it updates the metadata in *Database* to keep the name of the previous file.

It's deployed by a Docker container.

## Project composition

In the [`static`](static) there are the static images for the *Mega* services.

In the folder [`handlers`](handlers) there are the functions that handles the transformation of the images (update, delete and squared).

[`cmd/main.go`](cmd/main.go) is the main file that exposes the static files and the API Server. It validates the requests with *MegaAuth*.

## Dependencies

The dependencies and the scripts are defined in the file [`go.mod`](go.mod).

## Environment Variables

In: `.env.development` || `.env.production`

```javascript
APP_ENV=string
MONGODB_URI=string
MEGAGOAL_SERVER_API_URL=string
VALIDATE_URI=string
```

## Development server

Run `GO_ENV=development go run cmd/main.go` for a dev server.

## Build & Deploy

[`Dockerfile`](Dockerfile) file builds the app for production and generates de Docker container.

[`docker-compose.yml`](docker-compose.yml) file manages the image and handle it easily within the *Mega* network.