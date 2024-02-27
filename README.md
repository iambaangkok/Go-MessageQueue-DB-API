# Cloud-Course-Docker-Project
A minimal high-performance containerized microservice API built using Go, Kafka, and MySQL.


# Commands

## Docker build
```sh
# normal
docker build -t iambaangkok/go-kafka-gateway-api:latest -f ./Dockerfile-gateway-api .
docker build -t iambaangkok/go-kafka-internal-api:latest -f ./Dockerfile-internal-api .
# with logs
docker build -t iambaangkok/go-kafka-gateway-api:latest --progress=plain --no-cache -f ./Dockerfile-gateway-api .
docker build -t iambaangkok/go-kafka-internal-api:latest --progress=plain --no-cache -f ./Dockerfile-internal-api .
```

## Docker push
```sh
docker push iambaangkok/go-kafka-gateway-api:latest
docker push iambaangkok/go-kafka-internal-api:latest
```

## Run compose file
```sh
docker compose down
docker compose up -d
```

# Presentaion Slides
https://docs.google.com/presentation/d/1BP5vfeZNJPy1zRIon3dRVjPaCpdui4-M6QJM9CAhTn0/edit?usp=sharing