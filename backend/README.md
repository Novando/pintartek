# Pasuwado Backend

## How to run

First of all duplicate the `config.example.json` into `config.local.json` in `config` directory. Fill all the field

### Non-docker (systemd, pm2, daemon, etc)

To run the app on non-docker system, PostgresSQL and Redis should be installed

```shell
# Migrate database
DB_NAME=your_db_name DB_USER=your_db_user DB_HOST=your_db_host DB_PASS=your_db_password make migrate 

# Build the service
SERVICE=passvault-service make build

# Run the service
./bin/passvault-service
```

### Docker

```shell
# Build the image
IMAGE=pasuwado VERSION=latest make docker 

# Run the service
sudo docker run pasuwado
```

## Project Structure