FROM quay.io/nathansparks/sim-recon-deps:centos7
MAINTAINER Nathan Sparks <nsparks@jlab.org>
RUN /bin/bash -c "hdpm build sim-recon gluex_root_analysis \
    && hdpm clean --obliterate"
WORKDIR /home/gx
CMD source /home/gx/env-setup/master.sh && /bin/bash
