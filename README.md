# Ridesharing

A web app to organize shared rides.

## Build
```
git clone git@github.com:Mexx77/ridesharing.git
cd ridesharing
docker build -t ridesharing .
```

## Run
```
MONGO_PW=
JWT_SECRET=
docker run \
 -p 8090:8090 \
 -e ENVIRONMENT=prod \
 -e MONGO_PW=$MONGO_PW \
 -e JWT_SECRET=$JWT_SECRET \
 --name ridesharing \
 --rm \
 ridesharing
```