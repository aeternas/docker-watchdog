#!/bin/bash
docker stop swadeshness-nginx & docker rm swadeshness-nginx:master swadeshness-nginx:development & docker rmi aeternas/swadeshness-nginx:development aeternas/swadeshness-nginx:master & docker run -v swadeshness-certs:/etc/nginx/certs/ -p 8080:8080 -p 443:443 -d --restart=always --name swadeshness-nginx --network swadeshness aeternas/swadeshness-nginx
