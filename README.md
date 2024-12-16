# Octo-Manager

![Testing Check](https://github.com/jorgebaptista/octo-manager/actions/workflows/test.yml/badge.svg)
![Linting](https://github.com/jorgebaptista/octo-manager/actions/workflows/lint.yml/badge.svg)
![Security Check](https://github.com/jorgebaptista/octo-manager/actions/workflows/security.yml/badge.svg)

Octo-Manager is a REST API service written in Go to manage GitHub repositories and pull requests, with deployment on Minikube and a CI/CD pipeline for testing, linting, and security checks.

## Overview

Octo-Manager provides endpoints to:

- Create GitHub repositories.
- Delete GitHub repositories.
- List GitHub repositories.
- List open pull requests for a given repository.

## Setup Instructions

### Prerequisites

- Go 1.23 installed.
- GitHub Personal Access Token (PAT) with repo permissions.
- Docker Desktop (with WSL 2 integration if on Windows).
- Ubuntu (WSL Distribution).

Find instrunctions on how to set up Docker [here](https://docs.docker.com/desktop/).

### Setup

1. Clone the Repository:

```bash
git clone https://github.com/jorgebaptista/octo-manager.git
cd octo-manager
```

2. Install Dependencies:

```bash
go mod tidy
```

3. Create a `.env` file in the root directory with your Personal Access Token details.

```bash
touch .env
```

```plaintext
GITHUB_TOKEN=your-github-token
GITHUB_OWNER=your-github-username-or-org
```

Replace `your-github-token` and `your-github-username-or-org` with your actual GitHub Personal Access Token (PAT) and username/organization.

4. Build the Docker Image:

```bash
docker build -t octo-manager:latest .
```

5. Run the Docker container and pass the .env file

```bash
docker run --env-file .env -p 8080:8080 octo-manager:latest
```

The service will be available at `http://localhost:8080`

## API Endpoints

- Create a Repo:
`POST /repos`
Request Body (JSON): `{"name": "new-repo-name"}`

- Delete a Repo:
`DELETE /repos/:name`
Path parameter `:name` is the repository name.

- List All Repos:
`GET /repos`

- List Pull Requests:
`GET /repos/:name/pulls`
Path parameter `:name` is the repository name.
Optional query parameter `?n=x` to limit the number of PRs.
