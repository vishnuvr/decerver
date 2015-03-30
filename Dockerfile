# Dependencies

## Make sure your machine has >= 1 GB of RAM.

## Make sure go version >= 1.3.3 is installed and set up
FROM golang:1.4
MAINTAINER Eris Industries <contact@erisindustries.com>

### The base image kills /var/lib/apt/lists/*.
RUN apt-get update && apt-get install -qy \
  jq libgmp3-dev

# Setup base
RUN mkdir --parents $GOPATH/src/github.com/eris-ltd
WORKDIR $GOPATH/src/github.com/eris-ltd

# Install EPM
RUN git clone https://github.com/eris-ltd/epm-go
RUN git clone https://github.com/eris-ltd/eris-std-lib
RUN cd epm-go/cmd/epm; go get .

# Install Decerver
COPY . $GOPATH/src/github.com/eris-ltd/decerver
RUN cd decerver/cmd/decerver; go get .

# Configure
ENV user eris
RUN groupadd --system $user && useradd --system --create-home --gid $user $user
COPY docker/config /home/$user/.decerver/
COPY docker/config.json /home/$user/.decerver/languages/
RUN chown --recursive $user /home/$user/

# Expose ports.
## HTTP receiver
EXPOSE 3000 3005

## IPFS
EXPOSE 4001 5001 8080

## Thelonious
EXPOSE 15254 30303

# Hard-code 'eris' instead of using $user to support Docker versions < 1.3.
USER eris
CMD ["decerver"]