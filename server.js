'use strict';

var
  express = require('express'),
  request = require('request'),

  server;

server = express();

if (process.env.SOURCE) {
  console.log("Proxying IPFS objects from " + process.env.SOURCE + ".");

  server.get(/.*/, function (serverRequest, response) {
    request('http://localhost:8080/ipfs/' + process.env.SOURCE + '/'
      + serverRequest.url).pipe(response);
  });
}
else {
  console.log("Serving local files from 'local'.");
  server.use('/', express.static('local'));
}

server.listen(3000, function () {
  console.log("Decerver listening at port 3000.");
});
