FROM golang:1.4
MAINTAINER Eris Industries <support@ErisIndustries.com>

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get upgrade -qy && apt-get install -qy \
  curl \
  git \
  ruby

RUN gem install foreman

# Install Cloud Foundry Node.js buildpack.
# http://docs.cloudfoundry.org/buildpacks/node/index.html

RUN git clone https://github.com/cloudfoundry/nodejs-buildpack.git \
  /opt/buildpack

WORKDIR /opt/buildpack
RUN git checkout v1.1.1
RUN git submodule update --init --recursive

# Install IPFS (http://ipfs.io/docs/install/).
RUN go get -u github.com/jbenet/go-ipfs/cmd/ipfs
RUN ipfs init

# Choose the user id number 1000 to work well with Kitematic volumes.
RUN useradd --system --uid 1000 node

COPY cmd.sh /
CMD ["/cmd.sh"]
