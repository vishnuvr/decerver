# Eris Decerver

## Introduction

Eris decentralized applications are *front-end only* browser applications written in standard HTML, CSS, and JavaScript.  The source code for each application is loaded via [IPFS](http://ipfs.io/) and served via a web server at port 3000.

## Running Decerver

### Environment Variables

<table>
<tr><td>SOURCE</td><td>the IPFS hash of the source code directory</td></tr>
</table>

The following examples assume the source code for the application is in a directory named `source`.

### Production

Set the `SOURCE` environment variable to the IPFS hash of the source code directory.

	export SOURCE=$(ipfs add -recursive -quiet source | tail -n -1)

    docker run \
      --name=hello-world \
      --env SOURCE \
      --publish 3000:3000 \
      --detach \
      eris/decerver:browser
      
### Development

When you want to make rapid changes to the source code, you can also load it from a local volume instead of from IPFS by mapping the volume to `/usr/src/app/local`:

    docker run \
      --name=hello-world \
      --volume $PWD/source:/usr/src/app/local \
      --publish 3000:3000 \
      --detach \
      eris/decerver:browser

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