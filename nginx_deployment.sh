#!/bin/bash
imageName=$1:$2
containerName=autodeployed-$3-$2
PORTS_PAIR=""

if [[ $2 == "master" ]]; then
        PORTS_PAIR="-p 443:443 -p 8086:8086"
else
        PORTS_PAIR="-p 8080:8080 -p 8084:8084"
fi

docker stop $containerName && docker rm -f $containerName && docker rmi $imageName
docker run -v swadeshness-certs:/etc/nginx/certs/ $PORTS_PAIR --restart=always --name $containerName --network swadeshness $imageName
