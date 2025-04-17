# Docker Deployment

## Backend Deployment

```bash
# Build backend image
docker build -t i18n-flow-backend ./admin-backend

# Run backend container
docker run -d \
  -p 8080:8080 \
  --env-file .env \
  i18n-flow-backend
```

## Frontend Deployment

```bash
# Build frontend image
docker build -t i18n-flow-frontend ./admin-frontend

# Run frontend container
docker run -d \
  -p 80:80 \
  i18n-flow-frontend
```
