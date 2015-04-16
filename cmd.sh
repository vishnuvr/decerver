#!/bin/sh

# Launch the IPFS daemon.
ipfs init
ipfs config Addresses.API /ip4/0.0.0.0/tcp/5000
ipfs daemon -writable &

echo Waiting for the IPFS daemon to be ready.

while ! curl --silent localhost:5000; do
  sleep 1
done

npm start
