FROM nathansparks/sim-recon-deps:ubuntu16
MAINTAINER Nathan Sparks <nsparks@jlab.org>
RUN /bin/bash -c "hdpm build --all \
    && hdpm clean --obliterate --all"
WORKDIR /home/gx
CMD /bin/bash -c "source /home/gx/.hdpm/env/master.sh && /bin/bash"
