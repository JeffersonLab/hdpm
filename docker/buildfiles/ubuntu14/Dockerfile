FROM nathansparks/sim-recon-deps:ubuntu14
MAINTAINER Nathan Sparks <nsparks@jlab.org>
RUN /bin/bash -c "source setup.sh && julia src/hdpm.jl build sim-recon \
    gluex_root_analysis && julia src/distclean.jl yes"
WORKDIR /home/hdpm/docker/data
CMD /bin/bash -c "cd /home/hdpm && source pkgs/env-setup/master.sh && /bin/bash"
