#!/bin/sh

# Launch the IPFS daemon.
ipfs daemon -writable &

echo Waiting for the IPFS daemon to be ready.

while ! curl --silent localhost:5000; do
  sleep 1
done

npm start
