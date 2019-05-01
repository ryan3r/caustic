# Caustic

## Run the application
```
docker-compose up -d
./caustic-windows # Replace -window with your os
```

# Better build
```
DOCKER_BUILDKIT=1 docker build -t ryan3r/caustic caustic -f caustic/buildkit.Dockerfile
```