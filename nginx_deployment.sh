#!/bin/bash
docker stop swadeshness-nginx

docker rm swadeshness-nginx 
docker rmi aeternas/swadeshness-nginx 
docker run -v swadeshness-certs:/etc/nginx/certs/ -p 8080:8080 -p 443:443 -p 8084:8084 -d --restart=always --name swadeshness-nginx --network swadeshness aeternas/swadeshness-nginx
