# docker build -t sim-recon:ubuntu14 .
FROM nathansparks/sim-recon-deps:ubuntu14
MAINTAINER Nathan Sparks <nsparks@jlab.org>
RUN /bin/bash -c "source setup.sh && julia src/hdpm.jl build sim-recon"
WORKDIR /home/hdpm/docker/data
CMD cd /home/hdpm && source pkgs/env-setup/hdenv.sh && /bin/bash
