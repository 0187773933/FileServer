#!/bin/bash
APP_NAME="public-blogger-server"
sudo docker rm -f $APP_NAME || echo ""
id=$(sudo docker run -dit \
--name $APP_NAME \
--restart="always" \
--network=6105-buttons-1 \
--mount type=bind,source="$(pwd)"/config.yaml,target=/home/morphs/config.yaml \
-p 5155:5155 \
$APP_NAME /home/morphs/config.yaml)
sudo docker logs -f $id