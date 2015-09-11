# docker build -t sim-recon:ubuntu14 .
FROM nathansparks/sim-recon-deps:ubuntu14
MAINTAINER Nathan Sparks <nsparks@jlab.org>
RUN /bin/bash -c "source setup.sh && julia src/hdpm.jl build sim-recon"
WORKDIR /home/work/data
CMD cd ../hdpm && source pkgs/env-setup/hdenv.sh && /bin/bash
