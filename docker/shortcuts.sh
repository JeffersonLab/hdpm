# docker shortcuts
alias docker-ip="docker inspect --format {{.NetworkSettings.IPAddress}}"
docker-rm() { docker rm $(docker ps -aq); }
docker-rmi() { docker rmi $(docker images -f "dangling=true" -q); }
alias docker-run="docker run -it --rm"
WORKDIR=$(pwd)
dock() { docker run -it --rm -v $WORKDIR:/home/gx sim-recon \
    /bin/bash -c "source /home/gx/env-setup/master.sh &&
    export CCDB_USER=$USER && $1"; }
