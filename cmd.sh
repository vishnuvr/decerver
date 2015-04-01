#!/bin/sh

# Launch the IPFS daemon.
ipfs daemon &

echo Waiting for the IPFS daemon to be ready.

while ! curl --silent localhost:5001; do
  sleep 1
done

npm start
