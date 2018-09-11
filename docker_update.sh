#!/bin/bash
imageName=$1:$2
containerName=$1-$2

docker build -t $imageName -f Dockerfile  .

echo Stop and delete old container...
docker stop $containerName && docker rm -f $containerName && docker rmi $imageName

echo Run new container...
docker run -d --env-file env.list --restart=always -p 82:8082 --name $containerName $imageName
