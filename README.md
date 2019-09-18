# Ridesharing

A web app to organize shared rides.

[![Run on Google Cloud](https://storage.googleapis.com/cloudrun/button.svg)](https://console.cloud.google.com/cloudshell/editor?shellonly=true&cloudshell_image=gcr.io/cloudrun/button&cloudshell_git_repo=https://github.com/Mexx77/ridesharing.git)


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