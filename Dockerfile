FROM dockercore/golang-cross:1.12.1@sha256:8541e3aea7b2cffb7ac310af250e34551abe2ec180c77d5a81ae3d52a47ac779 AS build
ENV     DISABLE_WARN_OUTSIDE_CONTAINER=1

RUN apt-get install -y -q --no-install-recommends \
    coreutils \
    util-linux \
    uuid-runtime

WORKDIR /go/src/github.com/docker/cli

RUN git clone https://github.com/docker/cli . && git checkout 80918147ff32b0bf2a424d9133a45bab670793ff

RUN make cross binary && \
 cp build/docker-linux-amd64 /usr/bin/docker

WORKDIR /go/src/github.com/docker/app/

# main dev image
FROM build AS dev
ENV PATH=${PATH}:/go/src/github.com/docker/app/bin/
ARG DEP_VERSION=v0.5.0
RUN curl -o /usr/bin/dep -L https://github.com/golang/dep/releases/download/${DEP_VERSION}/dep-linux-amd64 && \
    chmod +x /usr/bin/dep
RUN go get -d gopkg.in/mjibson/esc.v0 && \
    cd /go/src/github.com/mjibson/esc && \
    go build -v -o /usr/bin/esc . && \
    rm -rf /go/src/* /go/pkg/* /go/bin/*
COPY . .

FROM dev AS cross
ARG EXPERIMENTAL="off"
ARG TAG="unknown"
RUN make EXPERIMENTAL=${EXPERIMENTAL} TAG=${TAG} cross

FROM cross AS e2e-cross
ARG EXPERIMENTAL="off"
ARG TAG="unknown"
# Run e2e tests
RUN make EXPERIMENTAL=${EXPERIMENTAL} TAG=${TAG} e2e-cross
