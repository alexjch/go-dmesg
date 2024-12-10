# go-dmesg
Golang dmesg stream and decoder library

## Usage
This is a utility library to read `/dev/kmesg`
```
import "github.com/alexjch/go-dmesg/pkg/dmesg"
```

## Quick start
Sample code to understand how to use the library can be run as a command.
```
git clone https://github.com/alexjch/go-dmesg.git
cd go-dmesg
go build cmd/go-dmesg.go
sudo ./go-dmesg
```

It's possible to inject a message for debugging with the following command:
```
sudo bash -c 'echo "Hello world" >/dev/kmsg'
```
