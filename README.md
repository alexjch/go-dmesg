# go-dmesg
Golang dmesg stream and decoder library

## Usage
This is a utility library to read `/dev/kmesg`
```
import "github.com/alexjch/go-dmesg/pkg/dmesg"
```

The function that creates the scanner should close it:
```go
scanner, err := dmesg.NewScanner()
if err != nil {
	log.Fatal(err)
}
defer scanner.Close()

decoder := dmesg.NewDecoder(scanner)
for decoder.Scan() {
	record := decoder.Record()
	_ = record
}
if err := decoder.Err(); err != nil {
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
