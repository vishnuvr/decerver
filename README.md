[![Stories in Ready](https://badge.waffle.io/eris-ltd/deCerver.png?label=ready&title=Ready)](https://waffle.io/eris-ltd/deCerver)

## What's this?

Generalized deCerver infrastructure for decentralized applications. The deCerver acts as the hub for DApps.

## Overview

The `decerver` package defines the main types and interfaces. The following generalized interfaces are provided:

- API server (acts as the frontend for the deCerver)
- decentralized database (eg. ethereum forked blockchain)
- decentralized filesystem (eg. IPFS)
- contracts integrator
- user notifier (for two-factor authentication)
- scripts runner (for running scheduled jobs)
- p2p communications layer (TBD)

For each of these packages, the deCerver acts as the hub for the entire system. The entire system operates in the manner shown below:

![deCerver architecture](docs/images/deCerver Structure.png)

In the image shown above the deCerver package is in Blue. As stated, it acts as the hub for the system.

Packages which *cannot interact* with users or systems outside of the context of the deCerver are:

* the eris fork of the ethereum package
* the legal integrator which controls real contract factories and changes to real contracts
* the decentralized file system
* the p2p communication system (TBD)

Packages which *can interact* with users or systems outside of the context of the deCerver are:

* the scripts runner which is made to predominantly make API calls to other systems and return data to the ethereum layer
* the notifier which is made to ensure that private keys of users are retained in user's control and also to ensure that legalities surrounding `offer` and `acceptance` for valid contractual arrangements are abided by

## Config

Every module type should also have a config type, and must have methods `WriteConfig(path)`, `ReadConfig(path)`, and `SetConfig(interface{})`. If you wan't to set non-default configs from the decerver level, we can do it with ReadConfig and SetConfig (and also use them for implementing flags. But flags at the decerver level. Though, each module should provide a standalone cli with a set of flags, and for testing.) We will have to talk more about standardizing configuration. For convenience, for now, config types should have json bindings, eg:

```go
type ChainConfig struct{
    Port int        `json:"port"`
    Mining bool     `json:"mining"`
    MaxPeers int    `json:"max_peers"`
    ConfigFile string `json:"config_file"`
    RootDir string  `json:"root_dir"`
    Name string     `json:"name"`
    LogFile string  `json:"log_file"`
}
```
This way, SetConfig can be passed json read in from a config file, or from main, or can be passed an initialized config struct itself. SetConfig boilerplate necessary.

