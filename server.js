'use strict';

var
  express = require('express'),
  http = require('http'),
  httpProxy = require('http-proxy'),
  request = require('request'),

  applicationServer,
  proxy, ipfsServer;


// application web server

applicationServer = express();

if (process.env.SOURCE) {
  console.log("Proxying IPFS objects from " + process.env.SOURCE + ".");

  applicationServer.get(/.*/, function (serverRequest, response) {
    request('http://localhost:8080/ipfs/' + process.env.SOURCE + '/'
      + serverRequest.url).pipe(response);
  });
}
else {
  console.log("Serving local files from 'local'.");
  applicationServer.use('/', express.static('local'));
}

applicationServer.listen(3000, function () {
  console.log("Decerver listening at port 3000.");
});


// IPFS proxy to allow for CORS
// IPFS should handle this already but it doesn't seem to be working.  See:
// https://github.com/ipfs/go-ipfs/issues/1017

proxy = httpProxy.createProxyServer({target: 'http://127.0.0.1:5000'});

ipfsServer = http.createServer(function (request, response) {
  response.setHeader('Access-Control-Allow-Origin', '*');
  response.setHeader('Access-Control-Allow-Headers', 'Content-Type');

  if (request.method !== 'OPTIONS') {
    // Work around restriction here:
    // https://github.com/ipfs/go-ipfs/blob/79360bbd32d8a0b9c5ab633b5a0461d9acd0f477/commands/http/handler.go#L58-L70
    // We should get a better understanding of the restriction instead of
    // working around it.
    delete request.headers.referer;

    proxy.web(request, response, {}, function (error) {
      console.error(error);
    });
  }
  else
    // Handle CORS preflight check ourselves because IPFS doesn't correctly.
    // See: https://github.com/ipfs/go-ipfs/issues/1049
    response.end();
});

ipfsServer.listen(5001, function () {
  console.log("IPFS proxy listening on port 5001.")
});
