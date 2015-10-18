# some bash aliases and functions to use with docker
alias docker-ip="docker inspect --format {{.NetworkSettings.IPAddress}}"
docker-rm() { docker rm $(docker ps -aq); }
#docker-rmi() { docker rmi $(docker images -q); }
alias docker-run="docker run -it --rm"
DATA_DIR=/path/to/data/directory
dock() { docker run -it -v $DATA_DIR:/home/hdpm/docker/data sim-recon \
    /bin/bash -c "source /home/hdpm/pkgs/env-setup/hdenv.sh && $1"; }
