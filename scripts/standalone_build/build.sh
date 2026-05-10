#!/usr/bin/env bash
set -euo pipefail

# Delete build folder if it exists
echo "Cleaning build directory..."
rm -rf build
mkdir -p "build/static"

# Run docker compose build
echo "Building containers..."
docker compose -f "docker-compose-build.yaml" up -d --build --remove-orphans 
docker compose -f "docker-compose-build.yaml" down 

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
