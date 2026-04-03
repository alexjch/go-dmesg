# go-dmesg
Golang dmesg stream and decoder library

## Usage
This is a utility library to read `/dev/kmsg`
```
import "github.com/alexjch/go-dmesg/pkg/dmesg"
```

The function that creates the scanner should close it:
```go
import (
	"context"
	"log"

	"github.com/alexjch/go-dmesg/pkg/dmesg"
)

ctx := context.Background()

scanner, err := dmesg.NewScanner()
if err != nil {
	log.Fatal(err)
}
defer scanner.Close()

decoder, err := dmesg.NewDecoder(scanner)
if err != nil {
	log.Fatal(err)
}

err = decoder.Follow(func(r *dmesg.Record) {
	log.Println(r.String())
}, ctx)
if err != nil {
	log.Fatal(err)
}
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
