# Caustic
[![coverage report](https://git.linux.iastate.edu/319spring19projectrepositories/ko/group-33/badges/master/coverage.svg)](https://git.linux.iastate.edu/319spring19projectrepositories/ko/group-33/commits/master)

[Install docker desktop](https://www.docker.com/products/docker-desktop) (windows professional and mac)
[Install docker toolbox](https://docs.docker.com/toolbox/toolbox_install_windows/) (for windows home)
Install docker on linux run `curl -sSL https://get.docker.com | sudo bash`

# Build and run

## Build the application with docker
```
docker-compose build
```

## Run the application
```
docker-compose up
```

## Build and run
```
docker-compose up --build
```

# Shutdown the server
```
docker-compose down
```

# Shutdown and clear mysql data
```
docker-compose down -v
```

# Better build
```
DOCKER_BUILDKIT=1 docker build -t ryan3r/caustic caustic -f caustic/buildkit.Dockerfile
```