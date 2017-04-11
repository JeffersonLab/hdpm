# docker shortcuts
alias docker-ip="docker inspect --format {{.NetworkSettings.IPAddress}}"
docker-rm() { docker rm $(docker ps -aq); }
docker-rmi() { docker rmi $(docker images -f "dangling=true" -q); }
alias docker-run="docker run -it --rm"
if [[ -z "$WORKDIR" ]]; then
    WORKDIR=$(pwd)
fi
dock() {
    local w=/home/gx/work
    docker run -it --rm -w $w -v $WORKDIR:$w sim-recon \
    /bin/bash -c ". /home/gx/.hdpm/env/master.sh &&
    export CCDB_USER=$USER && $1"
}
