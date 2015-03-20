# Eris Decerver

## Introduction

Eris decentralized applications are written in [Node.js](https://nodejs.org/)
and built with the [Cloud Foundry Node.js buildpack](http://docs.cloudfoundry.org/buildpacks/node/).  You may want to read Cloud Foundry's [Tips for Node.js Applications](http://docs.cloudfoundry.org/buildpacks/node/node-tips.html).

The source code for your application is loaded via [IPFS](http://ipfs.io/).

You need to use the `VCAP_APP_PORT` environment variable to determine which port
your application should listen on:

    app.listen(process.env.VCAP_APP_PORT);

## Running Decerver

### Source Code Format

The source code for the application is bundled with [GNU Tar](https://www.gnu.org/software/tar/).  This is because IPFS doesn't yet support symbolic links.

### Environment Variables

<table>
<tr><td>SOURCE</td><td>the IPFS hash of the source code archive</td></tr>
<tr><td>VCAP_APP_PORT</td><td>choose a non-system port</td></tr>
</table>

The following examples assume the source code for the application is in a directory named `source`.

### Production

Set the `SOURCE` environment variable to the IPFS hash of the source code archive.

	cd source
	  tar --create --file=/tmp/source.tar *
	cd ..
	
	export SOURCE=$(ipfs add -quiet /tmp/source.tar | tail -n -1)
	export VCAP_APP_PORT=3000

    docker run \
      --name=hello-world \
      --detach \
      --env SOURCE --env VCAP_APP_PORT \
      --publish 3000:$VCAP_APP_PORT \
      eris/decerver:node
      
### Development

When you want to make rapid changes to the source code, you can also load it from a local volume instead of from IPFS by mounting the volume `/home/node`:

    docker run \
      --name=hello-world \
      --detach \
      --env SOURCE --env VCAP_APP_PORT \
      --publish 3000:$VCAP_APP_PORT \
      --volume $PWD/source:/home/node \
      eris/decerver:node

You may want to launch your application with a program that watches for changes in the source code and then restarts the application automatically, like [nodemon](http://nodemon.io/).  You would do this by adding nodemon as a dependency in your `package.json` file with an appropriate `start` script:

```
  "dependencies": {
    "nodemon": "^1.3.7"
  },
  "scripts": {
    "start": "nodemon server.js",
  }
```

# Copyright

Copyright 2015 Eris Industries

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.