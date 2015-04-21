FROM node:0.10-onbuild
MAINTAINER Eris Industries <support@ErisIndustries.com>

# Install IPFS.
RUN npm install --global go-ipfs

## Allow access to the API from the browser.
ENV API_ORIGIN *

# Serve README by default.
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get upgrade -qy && apt-get install -qy markdown
RUN mkdir local
RUN markdown README.md > local/index.html

# application web server
EXPOSE 3000

# IPFS API
EXPOSE 5001

COPY cmd.sh /
CMD ["/cmd.sh"]
