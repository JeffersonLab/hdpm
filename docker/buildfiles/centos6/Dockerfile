FROM quay.io/nathansparks/sim-recon-deps:centos6
MAINTAINER Nathan Sparks <nsparks@jlab.org>
RUN scl enable python27 "source setup.sh && julia src/hdpm.jl build sim-recon"
WORKDIR /home/hdpm/docker/data
CMD cd /home/hdpm && source pkgs/env-setup/hdenv.sh && /bin/bash
