# rtorrent-cleaner

rtorrent-cleaner offer you a fast way to compare all files in your download directory and all torrent currently managed by your rtorrent instance.

rtorrent-cleaner use ruTorrent RPC API.

## Installation

- Download rtorrent-cleaner-X.Y.Z (where X.Y.Z are replace by the version number)
- Copy the rtorrent-cleaner-X.Y.Z anywhere in your path

## Getting Started

### Basic run 

```
rtorrent-cleaner -ruTorrent=https://rtorrent.mydomain.com/ -dir=/home/user/download
```

### Run disabling SSL/TLS check certificate (useful with self-signed certificate)

```
rtorrent-cleaner -ruTorrent=https://rtorrent.mydomain.com/ -dir=/home/user/download -no-check-certificate
```

### Run with absolute path for output

```
rtorrent-cleaner -ruTorrent=https://rtorrent.mydomain.com/ -dir=/home/user/download -absolute-path
```

### Display informations about rtorrent-cleaner

```
rtorrent-cleaner -version
```

## Compile

- clone repo
```
git clone https://github.com/kuroidoruido/rtorrent-cleaner.git
```

- download dependencies
```
make install
```

- compile
```
make compile
```

- edit RUN_ARGS variable at start of the Makefile to match your installation
- run
```
make run
```
