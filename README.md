# Setup

## Instal mise

Linux:

```curl https://mise.run | sh```

Windows:

```winget install jdx.mise```

## Activate mise

https://mise.en.dev/getting-started.html#activate-mise

## Install dependencies

```mise install```


# Local Development

```just dev```

# Local Deployment With Docker/Podman

```just dev-docker```

Then in browser navigate to http://localhost with admin/password credentials

Import entries and rules.

# Local build for standalone deployment

## Prerequisites

- Linux (or WSL)
- working docker/podman under user account
- installed zip

## Build

- Open Linux or WSL session
- Run ```just build```
- copy generated zip file to deployment servers

## Test

- have container runtime installed

```just test```

or, manually:

- systemctl --user start podman.socket
- export DOCKER_HOST=unix://$(podman info --format '{{.Host.RemoteSocket.Path}}')
- go test ./...

# Dependency updates

## Scan

Run:
```just scan```

Observe any error in the output

### Frontend

- navigate to frontend folder and run
    - npm outdated
    - npm update
    - npm run dev (and check for error)

### Backend

- based on the findings update mise to use required version of golang
- Update webapi/Dockerfile* to use the same golang image
- ```just update-backend```

Now that both frontend and backend are updated create new commit/release/build.