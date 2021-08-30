FROM scratch
MAINTAINER dev@gomicro.io

ADD train train

CMD ["/train"]
