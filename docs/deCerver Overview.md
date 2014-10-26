## Overview

The deCerver is the operational center for Eris Industries user facing stack. As such, it manages connections to: 

* a blockchain layer,
* a torrent layer,
* a decentralized file system layer,
* a secure peer-to-peer communications layer,
* a user facing web server layer, and
* a peer discovery mechanism. 

It should have interfaces for each of the above aspects which are agnostic to the actual implementations. Throughout the rest of this document the above layers will be referred to as "modules".

## Structure

The deCerver is the core library, should declare these interfaces, and should import nothing but core golang libraries. To actually implement a decerver, one writes a client which imports the deCerver library and the particular modules that implement the interfaces.

So for example, one might write the EFS library (file system interface), the Vortex library (torrent interface), a wrapper library to eth-go (implementing the blockchain interface), a wrapper to some DHT (peer discovery interface) and then write a client that imports all these and is used by a web structured DAPP (decentralized application). 

The client which is running on a users' computer, a personal cloud droplet, or on a centralized server (for corporate deployments) will initializes a deCerver object with the particular modules needed to run that DAPP. It will then call Decerver.Init, which will call the respective initializing routines for each of the objects required for the modules which are used by the DAPP. After modules have been initialized, the deCerver will start a web server, and serve the web structured DAPP to the user.

The advantage of this design is that it gives users of EI's deCerver (especially DAPP developers) a general purpose decerver which can swap in and out the various blockchains, torrent clients, file systems, peer discovery mechanisms, and and other modules as needed, so long as they are properly wrapped to implement the interfaces. 

Another advantages to this modular design is that we can wrap any golang library that already exists, so it implements the interface, and use those (e.g. using `eth-go`’s peer discovery, or `bitcoin`'s, a `torrent client`’s, `bitmessage`'s peer discovery, etc.). 

Something we might want to consider is replacing the `Vortex` TorrentClient with a generalized FileSharing module, which will mostly be implemented as torrent clients but could be anything. Also, we will probably want to add a peer-to-peer communication interface to the deCerver, but that will be developed sometime in early 2015.

## deCerver Configuration

The deCerver should have a configuration file which allows users to control access rules, allowed ports, allowed urls, TLS, etc. This configuration should be held in `~/.efs` (or an equivalent directory which we may implement over time). 

## deCerver Object and Struct

The deCerver object should look something like:

```go
type Decerver struct{
    Blockchain Blockchain
    TorrentClient TorrentClient
    FileSystem FileSystem
    PeerDiscover PeerDiscover
    WebApp WebApp
    Server Server
}
```

For the Server struct, we need to standardize some access rules. Let’s worry about that later. For now, a default server will be initialized, without TLS, looking like [CSK -> hmm. I'd like it to be fully https from ASAP if we can manage that. Reason being that one of the key business plans for EI is to host personal cloud-based instances of this stack for users (who will pay a nominal fee).]

```go
type Server struct{
    port int // listening port
    allowed_host []string // accepted url.Host (ie. localhost for now)
}
```

Eventually, we can add location of TLS certs, more allowed hosts, whatever.

## deCerver Interfaces

Now, the interfaces. I’m leaving `Init` routines without params for now. They will need though, obviously. Each of the modules should interface `GetInfo()` so the deCerver can know what implementation of the interface is being used, for logging purposes and possibly for some implementation specific customizations.

```go
type Blockchain interface{
    Init() 
    GetInfo()
    GetStorageAt(address, storage_addr) 
    Subscribe(address) // subscribe to updates at a particular address (reactors)
    Publish(...) // we’ll have to talk about this one. its more complex. 
}
```

```go
// replace with generalized FileSharer ?
type TorrentClient interface{
    Init() 
    GetInfo()
    CreateTorrent(files []string, name string)
    DownloadTorrent(btih string)
    CheckTorrent(btih string)
    Cleanup() // maybe?
}
```

```go
// we have to flesh this out some more
// yes we’ll use EFS, but let’s generalize the functionality so I can swap in standard unix or whatever i want. It shouldn’t have to know about commits, I don’t think, but maybe it should have a notion of blobs and trees? They do mark a rather generalized FS
type FileSystem interface{
    Init()
    GetInfo()
    CheckObj(name string) // does object exist?
    WriteBlob(contents []byte, name string)
    ReadBlob(name string)
    WriteTree(contents [][]byte, names []string)
    ReadTree(name string)
    CheckoutTree(name string) // checkout tree into a new working dir
}
```


```go
// needs to implement various routing functions for different urls
// also, needs to specify conversations with BTFP. This might be tricky
// maybe for now we build a standard webapp into the decerver
// and then we can work on generalizing
// this should be the only interface that might have to know about the others
// generalization might be tricky, but fun
type WebApp interface{
    Init()
    GetInfo()
    ...    
}
```

```go
// again, most useful for whisper. but this can wrap a DHT or the TorHiddenService idea for better privacy
type PeerDiscover interface{
    Init()
    GetInfo()
    FindPeer(value []byte) // find peer possesing the value information
    AnnouncePeer(ip []byte, port int, value []byte)
    RelayPeers( … ) // need to think about this more 
}
```

All of this will need to change and adapt as we go along, but at least this is a basis for how everything will fit together. Thoughts?

