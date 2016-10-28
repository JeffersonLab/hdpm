FROM quay.io/nathansparks/sim-recon-deps:centos6
MAINTAINER Nathan Sparks <nsparks@jlab.org>
RUN scl enable devtoolset-3 python27 "hdpm build sim-recon gluex_root_analysis \
    && hdpm clean --obliterate"
WORKDIR /home/gx
CMD source /home/gx/env-setup/master.sh && /bin/bash
