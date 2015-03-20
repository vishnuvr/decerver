#!/bin/sh

# Launch the IPFS daemon.
ipfs daemon &

# If there isn't already a node directory then we're in production and should
# grab the source from IPFS.
if [ ! -d /home/node ]; then
  echo Getting source code from $SOURCE.

  while ! curl localhost:5001; do
    sleep 1
  done

  mkdir /home/node
  cd /home/node
  ipfs cat $SOURCE | tar --extract
  chown --recursive node .
fi

# Build the Node.js application and start it.
su node << 'EOSU'
  cd /home/node
  /opt/buildpack/bin/detect .
  /opt/buildpack/bin/compile . /tmp

  for file in .profile.d/*.sh; do
    . "$file"
  done

  foreman start
EOSU
