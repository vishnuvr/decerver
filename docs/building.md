## Decerver Executable Build Sequence

### Step 1 - Get the `decerver` binary into path

#### If user has golang installed

```
go get -u -d github.com/eris-ltd/decerver
cd $GOPATH/src/github.com/eris-ltd/decerver/cmd/decerver
go install

go get -u -d github.com/eris-ltd/epm-go
cd $GOPATH/src/github.com/eris-ltd/epm-go/cmd/epm
go install
cd $GOPATH/src/github.com/eris-ltd/epm-go/cmd/iepm
go install

```

Will get the primary executable and dependencies.

```
sudo mv $GOPATH/bin/decerver /usr/bin/.
sudo mv $GOPATH/bin/epm /usr/bin/.
sudo mv $GOPATH/bin/iepm /usr/bin/.
```

#### If user does not have golang installed

Download latest binaries for platform from:

https://dl.erisindustries.com/decerver/latest
https://dl.erisindustries.com/epm/latest
https://dl.erisindustries.com/iepm/latest

#### Binary Population

We should use [goxc](https://github.com/laher/goxc) to populate the binaries.

### Step 2 - Pull Go Docs as Man Pages and Install to appropriate man location

@mattdf, you'll have to research how to most effectively do this.

-----

**Note** if only the backend is desired, stop here. To install the backend + the client continue

-----

### Step 3 - Install Atom-Shell

Use the atom-executable-example.sh in this directory to make an executable script and save to `/usr/local/bin`.

Download the latest atom-shell and extract to `/usr/local/share/decerver`. See an example of what is needed [here](https://github.com/atom/grunt-download-atom-shell)

### Step 4 - Install decerver-client

Pull latest decerver-client and save to `/usr/local/share/decerver/resources/app` per [this page](https://github.com/atom/atom-shell/blob/master/docs/tutorial/application-distribution.md#application-distribution)

### Step 5 - Execute decerver init command

### Step 6 - Download default DAPPs

* decerver-admin
* shittytube
* shittyplace

### Step 7 - Install .desktop file (if Linux) or Start Menu thingy (if Windows)