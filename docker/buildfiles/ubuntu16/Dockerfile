FROM hddeps:u16
RUN /bin/bash -c "hdpm install && hdpm clean --obliterate"
WORKDIR /home/gx
CMD /bin/bash -c "source /home/gx/.hdpm/env/master.sh && /bin/bash"
