# orochi

Like a distributed KVS

## How to make it work

### Build

```
# clone this repository and then:

$ cd orochi
$ go build
```

If the build succeeded, an executable binary `orochi` will be created on the directory .

### Launch orochi

```
# start an orochi with port 3000 (default)
$ ./orochi &

# start an orochi with specified port
$ ./orochi --port 3001 &
$ ./orochi --port 3002 &
```

Limitation: Currently port 3000, 3001 and 3002 can only be specified.

### Example to Get/Post key-value to orochi using curl


```bash
$ curl -X GET localhost:3000/hoge
# no value will be returned

$ curl -X POST localhost:3001/hoge  -d "fuga"
# hoge => "fuga" is stored

$ curl -X POST localhost:3002/hoge  -d "piyo"
# hoge => "piyo" is stored ("fuga" has been overwritten)

$ curl -X GET localhost:3000/hoge
piyo

$ curl -X GET localhost:3001/hoge
piyo

$ curl -X GET localhost:3002/hoge
piyo

$ curl -X POST localhost:3001/hoge -d "foobar"
# hoge => "foobar" is stored ("piyo" has been overwritten)

$ curl -X GET localhost:3000/hoge
foobar

$ curl -X GET localhost:3001/hoge
foobar

$ curl -X GET localhost:3002/hoge
foobar
```

## License

MIT

## Author

Yosuke Akatsuka (a.k.a pankona)
