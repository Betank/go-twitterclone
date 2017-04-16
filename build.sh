#!/bin/bash

for d in tweet-command-service tweet-query-service user-command-service user-query-service stats-service authentication-service; do
    cd $d

    echo -n "Building $d"
    CGO_ENABLED=0 GOOS=linux go build -o $d-docker
    
    echo -n "."
    docker build -t "twitterclone-$d" --quiet .

    cd ..
done

cd api-gateway/server

echo -n "Building api-gateway"
CGO_ENABLED=0 GOOS=linux go build -o api-gateway-docker
    
echo -n "Building api-gateway"
cd ..
npm install --quiet

echo -n "."
npm run build --quiet

echo -n "."
docker build -t "twitterclone-api-gateway" --quiet .