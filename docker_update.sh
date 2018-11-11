#!/bin/bash
imageName=$1:$2
containerName=autodeployed-$3-$2
PORTS_PAIR=""
ENV_FILE=""

if [[ $2 == "master" ]]; then
	PORTS_PAIR="80:8080"
	ENV_FILE="prod_env.list"
else
	PORTS_PAIR="8080:8080"
	ENV_FILE="env.list"
fi

docker build -t $imageName -f Dockerfile  .

echo Stop and delete old container...
docker stop $containerName && docker rm -f $containerName && docker rmi $imageName

echo Run new container...
docker run -v certs:/certs -d --env-file $ENV_FILE --restart=always -p $PORTS_PAIR --name $containerName $imageName
