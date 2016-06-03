FROM quay.io/nathansparks/sim-recon-deps:centos6
MAINTAINER Nathan Sparks <nsparks@jlab.org>
RUN scl enable devtoolset-3 python27 "source setup.sh && julia src/hdpm.jl build sim-recon \
    && julia src/distclean.jl yes"
WORKDIR /home/hdpm/docker/data
CMD cd /home/hdpm && source pkgs/env-setup/master.sh && /bin/bash
