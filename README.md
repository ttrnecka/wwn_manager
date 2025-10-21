# Local deployment with docker

Ensure docker is installed on the machine where are are planing to run this.

```
git clone https://github.com/ttrnecka/wwn_manager.git

cd wwn_manager

# update the .env file. set to IMPORT_ZSCALER_CERT=true if you run this behind zscaler

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