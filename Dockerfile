FROM ubuntu:14.04

RUN apt-get update && \
 apt-get -y install git make

RUN rm -rf /usr/local/go
ADD https://go.dev/dl/go1.18.4.linux-amd64.tar.gz /tmp/
WORKDIR "/tmp"
RUN tar -C /usr/local -xvf go1.18.4.linux-amd64.tar.gz

ARG CACHEBUST=1
RUN mkdir -p -v /root/.ssh/
RUN touch /root/.ssh/known_hosts
ADD data/deploy/id_ed25519 /root/.ssh/

WORKDIR "/root"
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

WORKDIR "/opt"
RUN git clone git@github.com:PlasmaTrout/swr.git
WORKDIR "/opt/swr"
ENV PATH="${PATH}:/usr/local/go/bin"
RUN make all

EXPOSE 5000

CMD /opt/swr/bin/server



