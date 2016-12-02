FROM nathansparks/sim-recon-deps:ubuntu16
MAINTAINER Nathan Sparks <nsparks@jlab.org>
RUN /bin/bash -c "hdpm build sim-recon gluex_root_analysis \
    && hdpm clean --obliterate --all"
WORKDIR /home/gx
CMD /bin/bash -c "source /home/gx/env-setup/master.sh && /bin/bash"
