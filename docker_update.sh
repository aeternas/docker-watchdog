#!/bin/bash
imageName=$1:$2
containerName=autodeployed-$2
PORTS_PAIR=""

if [[ $2 == "master" ]]; then
	PORTS_PAIR="80:8080"
else
	PORTS_PAIR="8080:8080"
fi

docker build -t $imageName -f Dockerfile  .

echo Stop and delete old container...
docker stop $containerName && docker rm -f $containerName && docker rmi $imageName

echo Run new container...
docker run -d --env-file env.list --restart=always -p $PORTS_PAIR --name $containerName $imageName
