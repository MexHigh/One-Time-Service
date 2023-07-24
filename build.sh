#!/bin/sh

docker build -t registry.git.leon.wtf/leon/one-time-service/amd64:latest . && docker push registry.git.leon.wtf/leon/one-time-service/amd64:latest