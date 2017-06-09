FROM quay.io/nathansparks/sim-recon-deps:centos6
MAINTAINER Nathan Sparks <nsparks@jlab.org>
RUN scl enable devtoolset-3 python27 "hdpm install \
    && hdpm clean --obliterate"
WORKDIR /home/gx
CMD source /home/gx/.hdpm/env/master.sh && /bin/bash
