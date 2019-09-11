#!/bin/bash
imageName=$1:$2
containerName=autodeployed-$3-$2
ENV_FILE=""

if [[ $2 != "master" ]] && [[ $2 != "development" ]]; then
  exit 0
fi

if [[ $2 == "master" ]]; then
	ENV_FILE="prod_env.list"
else
	ENV_FILE="env.list"
fi

docker build -t $imageName -f Dockerfile  .

echo Stop and delete old container...
docker stop $containerName && docker rm -f $containerName && docker rmi $imageName

echo Run new container...
docker run -v swadeshness-certs:/certs -d --env-file $ENV_FILE --restart=always --name $containerName --network swadeshness $imageName
