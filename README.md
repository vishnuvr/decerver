[![Stories in Ready](https://badge.waffle.io/eris-ltd/deCerver.png?label=ready&title=Ready)](https://waffle.io/eris-ltd/deCerver)[![GoDoc](https://godoc.org/github.com/decerver?status.png)](https://godoc.org/github.com/eris-ltd/decerver)

![deCerver logo](docs/images/decerver-color.png)

## The Decerver

A worker for a decentralized stack. The Decerver is the application server for Distributed Applications (DAPPs).

## What is This Thing?

The Decerver helps developers build applications which leverage opt-in data ownership and significantly increased data utility for both customers and businesses – a software design paradigm we call **Participatory Architecture**. Using peer-to-peer and distributed systems, the Decerver allows the creation of web style, data-driven, interactive distributed applications that can be safely, securely, and reliably deployed and managed. The Decerver significantly lowers the barrier to entry for the production, distribution and maintenance of distributed applications. All while allowing users to participate in the scaling and data security of the application.

More specifically, the Decerver is a distributed application server harmonizes actions across various modules which act as distributed file stores, distributed data stores, or other utility modules. The Decerver integrates distributed data stores (blockchains), a distributed filesystem, a scripting layer, and a legal integrator to incorporate [Legal Markdown](https://lmd.io/)-based contracts into the smart contract stack – effectively putting the “contract” into “smart contracts.”

Applications built for the Decerver are based on web design modalities. In other words, the user interfaces for these applications are written in languages which almost any developer, and even some [heads of state](http://techcrunch.com/2014/12/08/barack-obama-becomes-the-first-president-to-write-code/), can write. Applications for the Decerver use HTML, CSS, and Javascript to provide their user interface.

Each of the modules which the Decerver utilizes has an established interface which exposes functions to a javascript runtime that executes inside of the Decerver’s core. These exposed functions allow a distributed application developer to design and implement a distributed application almost entirely in javascript – with the exception that if a blockchain that uses smart contracts is needed for that distributed application, the developer will need to use one of the smart contract languages to build those smart contracts (and the other exceptions of html and css of course).

## Installation

You must have [Go](https://golang.org/) installed.

```
go get github.com/eris-ltd/decerver
cd $GOPATH/src/github.com/eris-ltd/decerver/cmd/decerver
go get -d .
go install
```

That's it! If you have problems building please do let us know by filing an issue here on Github. We will do our best to assist.

**Please note** at this time we have not effectively tested Decerver on Windows so if you have a windows machine we welcome your feedback if you run into any problems (or if you do not!).

## Usage

For Usage and Tutorials, please see the [Decerver](https://decerver.io) site.

## Contributions

1. Fork
2. Hack
3. Pull Request

Please note that any pull requests which are merged into this repository explicitly accept the licensing thereof.
