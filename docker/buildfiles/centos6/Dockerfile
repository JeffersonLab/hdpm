FROM hddeps:c6
RUN scl enable devtoolset-3 python27 "hdpm install \
    && hdpm clean --obliterate"
WORKDIR /home/gx
CMD source /home/gx/.hdpm/env/master.sh && /bin/bash
