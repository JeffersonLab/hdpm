# some bash aliases and functions to use with docker
alias docker-ip="docker inspect --format {{.NetworkSettings.IPAddress}}"
docker-rm() { docker rm $(docker ps -aq); }
#docker-rmi() { docker rmi $(docker images -q); }
alias docker-run="docker run -it --rm"
DATA_DIR=/path/to/data/directory
CCDB_USER=$USER
dock() { docker run -it -v $DATA_DIR:/home/hdpm/docker/data sim-recon \
    /bin/sh -c "source /home/hdpm/pkgs/env-setup/hdenv.sh &&
    export CCDB_USER=$CCDB_USER && $1"; }
