FROM hddeps:c7
RUN /bin/bash -c "hdpm install && hdpm clean --obliterate"
WORKDIR /home/gx
CMD source /home/gx/.hdpm/env/master.sh && /bin/bash
