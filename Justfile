build:
    chmod 700 build.sh
    ./build.sh

dev: dev-mongo
    tmux new-session -d 'cd frontend && npm run dev' \; \
         split-window 'cd webapi && go run .' \; \
         attach

dev-mongo: docker-stop
    docker compose -f docker-compose.yaml up -d --build --remove-orphans mongodb

docker: docker-stop
    touch .env && docker compose -f docker-compose.yaml up -d --build --remove-orphans 

docker-ls:
    docker compose -f docker-compose.yaml ps
    
docker-stop:
    docker compose -f docker-compose.yaml down

scan:
    docker compose -f docker-compose-scan.yaml up --build
    
test:
    #!/bin/bash
    set -e
    
    # Determine runtime
    if command -v podman &> /dev/null; then
        RUNTIME="podman"
    else
        RUNTIME="docker"
    fi
    echo "Using container runtime: $RUNTIME"
    
    # Start the appropriate socket and get the path
    if [ "$RUNTIME" = "podman" ]; then
        systemctl --user start podman.socket
        SOCKET_PATH="unix://$(podman info --format '{{{{.Host.RemoteSocket.Path}}')"
    else
        SOCKET_PATH="$(docker context inspect --format '{{{{.Endpoints.docker.Host}}')"
    fi
    
    # Give it a moment to start
    sleep 1
    
    echo "Socket path: $SOCKET_PATH"
    
    # Run tests with the appropriate Docker host
    DOCKER_HOST="$SOCKET_PATH" go test ./webapi/...

update-backend: 
    cd webapi && go get -u
