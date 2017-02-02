#!/bin/sh
sudo docker build -t go-images-server .
sudo docker run -p 9090:9090 -d -it --rm --name go-images-server go-images-server

