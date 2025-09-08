# Tutup Lapak Project Sprint – Microservices Scaffold

This repository scaffolds a microservices setup with a gateway and two services. The gateway (`backend-infra`) is built with Go Fiber, `auth-service` is Go Fiber, and `profile-service` is Nest.js.

## Layout

- `backend-infra` – Go Fiber gateway routing `/auth/*` and `/profile/*`.
  - `backend-infra/docker-compose.yml` – Local infra (Postgres + Redis) for development.
- `auth-service` – Go Fiber service (port 3001) with health endpoint.
- `profile-service` – Nest.js service (port 3002) with health endpoint.
- `.devcontainer/` – VS Code Dev Containers configuration attaching to the `dev` service.

## Quick start

1. Infra (optional, local only):
   - From `backend-infra/`, run: `docker compose up -d`
   - Postgres: `localhost:5432` (postgres/postgres)
   - Redis: `localhost:6379`

2. Open in Dev Container (VS Code) [optional]:
   - Install the Dev Containers extension.
   - “Reopen in Container”.

## Environment

- Ports by convention: Gateway 3000, Auth 3001, Profile 3002

## Notes

- The gateway strips the `/auth` and `/profile` prefixes before proxying to each service.
- The `dev` container mounts the repo at `/workspaces` and caches Go modules in a named volume.
