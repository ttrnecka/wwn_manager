#!/usr/bin/env bash
set -euo pipefail

# Delete build folder if it exists
echo "Cleaning build directory..."
rm -rf build
mkdir -p "build/static"

cleanup_container() {
  local container_name="$1"

  # Check if a container with this name/ID exists (running or stopped)
  if [ $(docker ps -aq -f name="^/${container_name}$") ]; then
    echo "Container '$container_name' exists. Removing..."
    # -f forces removal if the container is currently running
    docker rm -f "$container_name" >/dev/null
  fi
}
cleanup_container wwn_manager_frontend
cleanup_container wwn_manager_backend

# frontend build
echo "Building frontend..."
docker build ../../ -f ./frontend/Dockerfile -t wwn_manager_frontend_build
docker create --name "wwn_manager_frontend" "wwn_manager_frontend_build"
docker cp "wwn_manager_frontend:/app/dist/." "build/static/"

# backend build
echo "Building backend..."
docker build ../../ -f ./webapi/Dockerfile.build_windows -t wwn_manager_backend_build
docker create --name "wwn_manager_backend" "wwn_manager_backend_build"
docker cp "wwn_manager_backend:/app/wwn_manager.exe" "build/"

# Move .env.template into build directory
echo "Moving .env.template to build/..."
cp "../../.env.template" "./build/.env.template"

DATE=$(date +%Y%m%d)
ZIPFILE="wwn_manager_${DATE}.zip"

if [[ -f "${ZIPFILE}" ]]; then
  echo "Removing existing ${ZIPFILE}..."
  rm -f "${ZIPFILE}"
fi

echo "Zipping contents of build directory into ${ZIPFILE}..."
(
  cd "build"
  zip -r "../${ZIPFILE}" .
)

echo "✅ Done! Created ${ZIPFILE}"
