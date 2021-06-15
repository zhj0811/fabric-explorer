# Copyright peerfintech All Rights Reserved
#

FROM frolvlad/alpine-glibc:latest
ADD bin/apiserver  /opt/bin/
WORKDIR /opt/bin
CMD ["./apiserver start"]
