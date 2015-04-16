FROM golang:1.4
MAINTAINER Eris Industries <support@ErisIndustries.com>


# Install Node.js
# https://github.com/joyent/docker-node/blob/1a414011089f16390800995f469f5f08446baf7f/0.10/Dockerfile

# verify gpg and sha256: http://nodejs.org/dist/v0.10.31/SHASUMS256.txt.asc
# gpg: aka "Timothy J Fontaine (Work) <tj.fontaine@joyent.com>"
# gpg: aka "Julien Gilli <jgilli@fastmail.fm>"
RUN gpg --keyserver pool.sks-keyservers.net --recv-keys 7937DFD2AB06298B2293C3187D33FF9D0246406D 114F43EE0176B71C7BC219DD50A3051F888C628D

ENV NODE_VERSION 0.10.38
ENV NPM_VERSION 2.7.3

RUN curl -SLO "http://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.gz" \
  && curl -SLO "http://nodejs.org/dist/v$NODE_VERSION/SHASUMS256.txt.asc" \
  && gpg --verify SHASUMS256.txt.asc \
  && grep " node-v$NODE_VERSION-linux-x64.tar.gz\$" SHASUMS256.txt.asc | sha256sum -c - \
  && tar -xzf "node-v$NODE_VERSION-linux-x64.tar.gz" -C /usr/local --strip-components=1 \
  && rm "node-v$NODE_VERSION-linux-x64.tar.gz" SHASUMS256.txt.asc \
  && npm install -g npm@"$NPM_VERSION" \
  && npm cache clear

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app


# Install IPFS (http://ipfs.io/docs/install/).
RUN go get -u github.com/jbenet/go-ipfs/cmd/ipfs

## Allow access to the API from the browser.
ENV API_ORIGIN *

# Node.js -onbuild
COPY . /usr/src/app
RUN npm install

# application web server
EXPOSE 3000

# IPFS API
EXPOSE 5001

COPY cmd.sh /
CMD ["/cmd.sh"]
