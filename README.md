# go-extid

It can be valuable to internally use a serial integer as an ID without revealing that ID to the outside world. go-extid
uses AES-128 to convert to and from an external ID that cannot feasibly be decoded without the secret key.

This prevents outsiders from quantifying the usage of your application by observing the rate of increase of IDs as well
as provides protection against brute force crawling of all resources.

## Example Usage

```go
	prefix := "user"
	key := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	et, err := extid.NewType(prefix, key)
	if err != nil {
		return err
	}

  et.Encode(1) // => "user_13189a6ae4ab07ae70a3aabd30be99de"
```

## Performance

```
$ go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/jackc/go-extid
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkEncode-16    	 8812752	       126.3 ns/op	     112 B/op	       4 allocs/op
BenchmarkDecode-16    	14138239	        78.88 ns/op	      32 B/op	       2 allocs/op
PASS
ok  	github.com/jackc/go-extid	5.298s
```

## Other Implementations

* [PostgreSQL](https://github.com/jackc/pg-extid)
