# Caustic

[Install docker desktop](https://www.docker.com/products/docker-desktop) (windows professional and mac)
[Install docker toolbox](https://docs.docker.com/toolbox/toolbox_install_windows/) (for windows home)
Install docker on linux run `curl -sSL https://get.docker.com | sudo bash`

# Build and run

The first time you run mysql (aka db) it will need initialize the database. So run `docker-compose up db` and wait until you see `ready for connections` in the middle of a line or the text stops flying by.  

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
