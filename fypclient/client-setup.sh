#!/bin/bash
#Get repo absaloute location for mounting into the container
local_workdir=$(cd $(dirname $(dirname "${BASH_SOURCE[0]}")) >/dev/null 2>&1 && pwd)
main(){
    #Working directory in the container
    local container_workdir=/go/src/github.com/wade-sam/fypclient
    local host_directory=$(eval echo ~$USER)
    homedir=$(getent passwd "$USER" | cut -d: -f6)
    #echo $homedir
    #Identifying container name
    local container_name=fyp-client
    docker run --rm -it --name $container_name --volume \
        $local_workdir:$container_workdir \
        --volume $homedir:/backup \
        --workdir $container_workdir \
        golang:1.15.5-alpine3.12
    
    echo $(go build)

}

main