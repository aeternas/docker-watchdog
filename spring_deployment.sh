#!/bin/bash
imageName=$1:$2
containerName=autodeployed-$3-$2
ENV_FILE=""

echo Stop and delete old container...
docker stop $containerName && docker rm -f $containerName && docker rmi $imageName

docker run -d --network swadeshness --restart=always --name $containerName $imageName
cat ~/rsats4 | docker exec -i $containerName sh -c 'mkdir /root/.ssh && cat > /root/.ssh/id_rsa'
