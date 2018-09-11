#!/bin/bash
imageName=$1:$2
containerName=swadeshness-$2

docker build -t $imageName -f Dockerfile  .

echo Stop and delete old container...
docker stop $containerName && docker rm -f $containerName && docker rmi $imageName

echo Run new container...
docker run -d --env-file env.list --restart=always -p 8083:8083 --name $containerName $imageName
