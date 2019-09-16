# Ridesharing

A web app to organize shared rides.

## Build
```
git clone git@github.com:Mexx77/ridesharing.git
cd ridesharing
docker build -t ridesharing .
```

## Run locally without Kensis putting
```
docker run \
 -p 8090:8090 \
 -e ENVIRONMENT=prod \
 -v `pwd`/server/sqlite.db:/server/sqlite.db \
 --name ridesharing \
 --rm \
 ridesharing
```