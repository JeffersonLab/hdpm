FROM quay.io/nathansparks/sim-recon-deps:centos7
MAINTAINER Nathan Sparks <nsparks@jlab.org>
RUN /bin/bash -c "hdpm build --all \
    && hdpm clean --obliterate --all"
WORKDIR /home/gx
CMD source /home/gx/.hdpm/env/master.sh && /bin/bash
