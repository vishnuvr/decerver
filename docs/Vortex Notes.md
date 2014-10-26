##Vortex

Hi everyone, and Ethan in particular. These are some current issues with the vortex system.

Currently, the interface (mockup) in the decerver spec is this:

```go
type TorrentClient interface{
    Init() 
    GetInfo()
    CreateTorrent(files []string, name string)
    DownloadTorrent(btih string)
    CheckTorrent(btih string)
    Cleanup() // maybe?
}
```

###Functions

####Init

If the dht will use  ethereum address as dht node ide, maybe it should be passed through here.

####GetInfo

For "global settings", Taipei uses a "flags" object to store a bunch of settings for torrents, but those are based on command line arguments. The way I'm handling this is to keep a settings struct, 'VortexSettings', that keeps default values that can be modified through functions defined on it. These values are both the taipei torrent flags, but can also have decerver torrent management options later (there are none atm). The flags are passed automatically to new torrent sessons upon creation.

####CreateTorrent

I'm guessing this will use your special function, so I'll leave it empty.

####DownloadTorrent

This now returns a 'chan string' that is written to when the torrent is done.

Not sure if you should be able to pass an array of hashes here also, and get an array of channels back? I could just make a single btih version that wraps the hash in an array, calls the array-taker, and return a single channel.

####CheckTorrent

I added a struct that contains torrent data before we all got together, so that it could be sent to a web based GUI:


```go
// Torrent data that is formatted and ready for display.
type TData struct {
	Name 		string // ts.torrentFile
	Hs			string // ts.M.InfoHash
	Pieces		int // ts.goodPieces
	PiecesTot 	int // ts.totalPieces
	Size		int64 // ts.totalSize
	PeerNr		int // len(ts.peers)
	Downloaded	uint64 // ts.si.Downloaded
	Uploaded	uint64 // ts.si.Uploaded
	Left		uint64 // ts.si.Left
	Seeders		uint // ts.ti.Incomplete
	Leechers	uint // ts.ti.Complete
}
```

Too much? Too little? Either way it's available as a choice.

####Cleanup

Seems good to me. This can be used to stop the channel reading, and gracefully shut down torrent sessions etc., and be called from the decerver shutdown function. There is already a 'quit channel' copy-pasted from Taipei-Torrent's torrentloop so np.

###Overall structure

If it uses an Init function, but no factory method, does this mean it should be lazy loaded upon initing? I'm assuming then it will be a singleton object, and that Init will return it?

Also, the settings struct will have methods to set the flags (use upnp, enable local peer discovery etc.), but I guess just have 'GetInfo' return a pointer to settings and write the getters/setters using the standard Go rules for names.

Finally, I'm still testing this while writing, so there might have to be small changes made that I'm not yet aware of.

// Andreas
