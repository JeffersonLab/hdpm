FROM nathansparks/sim-recon-deps:ubuntu16
MAINTAINER Nathan Sparks <nsparks@jlab.org>
RUN /bin/bash -c "source setup.sh && julia src/hdpm.jl build sim-recon \
    && julia src/distclean.jl yes"
WORKDIR /home/hdpm/docker/data
CMD /bin/bash -c "cd /home/hdpm && source pkgs/env-setup/master.sh && /bin/bash"
