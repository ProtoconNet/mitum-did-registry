### mitum-did

*mitum-did* is a did contract model based on the second version of mitum(aka [mitum2](https://github.com/ProtoconNet/mitum2)).

#### Installation

Before you build `mitum-did`, make sure to run `docker run` for digest api.

```sh
$ git clone https://github.com/ProtoconNet/mitum-did

$ cd mitum-did

$ go build -o ./mitum-did
```

#### Run

```sh
$ ./mitum-did init --design=<config file> <genesis file>

$ ./mitum-did run <config file> --dev.allow-consensus
```

[standalong.yml](standalone.yml) is a sample of `config file`.

[genesis-design.yml](genesis-design.yml) is a sample of `genesis design file`.
