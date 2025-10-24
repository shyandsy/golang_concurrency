# golang_concurrency
enhance programming skills by handling high concurrency by golang

## case 1: 5-Seconds Countdown
This sample code demonstrates a 5-second countdown timer in go. It will display numbers from 5 to 1, followed by "finish". Ideal for learning basic concurrency and timeing in go

output
```
5
4
3
2
1
finished!

End Processing!!
```

## case 2: chan usage

1. create chan
- read/write chan
- read only chan
- write only chan

2. read & write chan
3. close chan
4. test
- try to write data into closed chan
- try to read from closed chan
5. buffered chan with quit signal chan

output

```shell
1 :  try_non_buffered_chan 
======================
inpu chan 1
read value from non buffered chan:  1 true
--------------------------------------------------------------------------

2 :  try_write_to_closed_chan 
======================
chan closed
will be panic next line..
Recovered from panic, reason:  send on closed channel
--------------------------------------------------------------------------

3 :  try_read_only_chan 
======================
read value from read only chan:  42 true
--------------------------------------------------------------------------

4 :  try_write_only_chan 
======================
read value from chan:  42
--------------------------------------------------------------------------

5 :  try_buffered_chan 
======================
worker:  2  waiting for assignment
worker:  0  waiting for assignment
worker:  1  waiting for assignment
assignment  2  finished by worker  1
worker:  1  waiting for assignment
assignment  0  finished by worker  2
worker:  2  waiting for assignment
assignment  1  finished by worker  0
worker:  0  waiting for assignment
assignment  5  finished by worker  0
worker:  0  waiting for assignment
assignment  4  finished by worker  2
worker:  2  waiting for assignment
assignment  3  finished by worker  1
worker:  1  waiting for assignment
assignment  8  finished by worker  1
worker:  1  waiting for assignment
assignment  7  finished by worker  2
worker:  2  waiting for assignment
assignment  6  finished by worker  0
worker:  0  waiting for assignment
assignment  9  finished by worker  1
worker:  1  waiting for assignment
workloads:  map[0:3 1:4 2:3]
```

## case 3: concurrent cache implemented by map

implement a concurrent cache based on map, usage for `sync.RWMutex`

auto expiration policy
- auto deletion by timer
- lazy deletion: check deadline on `Get` method

interface
```go
type Cache interface {
	Set(key string, value string, ttl time.Duration) // set ket value with ttl
	Get(key string) (string, bool) // get key, return "", false if not exist
	Close() // close cache
}
```

output
```shell
$ cd case3_mapcache

$ go test -v ./
=== RUN   TestSetGet
2025/10/24 15:23:33 autoExpiration quit
--- PASS: TestSetGet (0.00s)
=== RUN   TestExpiration
--- PASS: TestExpiration (3.00s)
=== RUN   TestPartialExpiration
2025/10/24 15:23:36 autoExpiration quit
--- PASS: TestPartialExpiration (10.00s)
PASS
ok      golang_concurrency/case3_mapcache       13.004s
```
