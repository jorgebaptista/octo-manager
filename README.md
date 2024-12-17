# Octo-Manager

[![CI/CD Pipeline](https://github.com/jorgebaptista/octo-manager/actions/workflows/pipeline.yaml/badge.svg)](https://github.com/jorgebaptista/octo-manager/actions/workflows/pipeline.yaml)

Octo-Manager is a REST API service written in Go that allows users to create, destroy, and list GitHub repositories, as well as list open pull requests for a specific repository.

## Features

Octo-Manager provides endpoints to:

- Create GitHub repositories.
- Delete GitHub repositories.
- List GitHub repositories.
- List open pull requests for a given repository with an optional limit.

## Setup Instructions

### Prerequisites

- Go 1.23
- GitHub Token with repo permissions
- Docker
- Minikube
- Kubectl

### Environment Variables

Create a `.env` file in the project root:

```plaintext
GITHUB_TOKEN=your_github_token
GITHUB_OWNER=your_github_owner
```

## Running Locally

1. **Build the Docker Image:**

```bash
docker build -t jorgebaptista/octo-manager:latest .
```

2. **Run the Docker Container:**

```bash
docker login
docker push jorgebaptista/octo-manager:latest
```

3. **Access the API:** Open `http://localhost:8080` in your browser.

## Deploying to Kubernetes (Minikube)

1. **Push the Docker Image to Docker Hub:**

```bash
docker login
docker push jorgebaptista/octo-manager:latest
```

2. **Apply Kubernetes Manifests:**

```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

3. **Access the Service:**

```bash
minikube service octo-manager
```

## API Endpoints

- **Create a Repo:** `POST /repos`

Request Body (JSON): `{"name": "new-repo-name"}`

- **Delete a Repo:**
`DELETE /repos/:name`

Path parameter `:name` is the repository name.

- **List All Repos:**
`GET /repos`

- **List Pull Requests:**
`GET /repos/:name/pulls`

Path parameter `:name` is the repository name.
Optional query parameter `?n=x` to limit the number of PRs.

## Testing

### Unit Tests

Run unit tests for the GitHub client:

```bash
go test ./tests/unit -v
```

### Integration Tests

Run integration tests for the API endpoints:

```bash
go test ./tests/integration -v
```

## CI/CD Pipeline

GitHub Actions workflows are set up to:

1. **Lint**: Run `golangci-lint`.
2. **Security**: Run `gosec`.
3. **Unit Tests**: Execute unit tests.
4. **Integration** Tests: Execute integration tests.
5. **Build & Push**: Build the Docker image and push it to Docker Hub.
6. **Deploy**: Deploy the latest Docker image to Kubernetes/Minikube.
