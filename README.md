# Local deployment with docker

Ensure docker is installed on the machine where are are planing to run this.

```
git clone https://github.com/ttrnecka/wwn_manager.git

cd wwn_manager

# if customization is needed create .env from .env.template file
# if not just run below

touch .env
docker compose -f docker-compose-dev.yaml up -d
```

Then in browser navigate to http://localhost with admin/password credentials

Import entries and rules.

# Local build for standalone deployment

## Prerequisites

- Linux (or WSL)
- working docker/podman under user account
- installed zip

## Build

- Open Linux or WSL session
- Run ./build.sh
- copy generated zip file to deployment servers

# Dependency updates

To update the dependencies you need local installation of Node.js and Go.

## Frontend

- Install Node v22.18 (or higher) on your server
- navigate to frontend folder and run
    - npm outdated
    - npm update
    - npm run dev (and check for error)

## Backend

- Install Go 1.24.x
- navigate to webapi and run
    - go get -u


Now that both frontend and backend are updated create new commit/release/build.