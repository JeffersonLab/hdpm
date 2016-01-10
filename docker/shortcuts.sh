# some bash aliases and functions to use with docker
alias docker-ip="docker inspect --format {{.NetworkSettings.IPAddress}}"
docker-rm() { docker rm $(docker ps -aq); }
#docker-rmi() { docker rmi $(docker images -q); }
alias docker-run="docker run -it --rm"
DATA_DIR=`pwd`/data
dock() { docker run -it --rm -v $DATA_DIR:/home/hdpm/docker/data sim-recon \
    /bin/bash -c "source /home/hdpm/pkgs/env-setup/hdenv.sh &&
    export CCDB_USER=$USER && $1"; }
