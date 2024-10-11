### mitum-did-registry

*mitum-did-registry* is a did contract model based on the second version of mitum(aka [mitum2](https://github.com/ProtoconNet/mitum2)).

#### Installation

Before you build `mitum-did-registry`, make sure to run `docker run` for digest api.

```sh
$ git clone https://github.com/ProtoconNet/mitum-did-registry

$ cd mitum-did-registry

$ go build -o ./mitum-did-registry
```

#### Run

```sh
$ ./mitum-did-registry init --design=<config file> <genesis file>

$ ./mitum-did-registry run <config file> --dev.allow-consensus
```

[standalong.yml](standalone.yml) is a sample of `config file`.

[genesis-design.yml](genesis-design.yml) is a sample of `genesis design file`.
