<p align="center">
  <h1 align="center">
    sesameFS
    <br/>
    <a href="https://github.com/godano/cardano-lib/blob/master/LICENSE" ><img alt="license" src="https://img.shields.io/badge/license-MIT%20License%202.0-E91E63.svg?style=flat-square" /></a>
  </h1>
</p>


The `cardano-node` team promised to work on a hot-loading mechanism that would make this project obselete. However, it hasn't been worked on so far, which is why we provide this userspace filesystem.

This project implements an userspace filesystem for Linux and FreeBSD to protect the operational certificate over most of its lifespan by providing a mechanism to store encrypted versions (pads) of the certificate files in a vault with a number of throw-away keys. The filesystem can be mounted to a specific directory (`key` in the example below) and it always contains exclusively the three files needed to operate the block producer.

```
key/
├── kes.skey
├── node.cert
└── vrf.skey
```

The filesystem provides an access control which would allow an one-time access to the three files given a valid throw-away key (with which the certificate has been encrypted) posted to the filesystem over a HTTP+JSON API. A user can decide how many keys they want to generate. As the name suggest, the key must be thrown away after usage and also the data encrypted with this key is thrown away. A detailed sequence diagram is shown in the *Concept* section.

<p align="center">
  <img src="https://upload.wikimedia.org/wikipedia/commons/1/1a/%D7%A2%D7%9C%D7%99_%D7%91%D7%90%D7%91%D7%90_%D7%9E%D7%AA%D7%97%D7%91%D7%90_%D7%A2%D7%9C_%D7%94%D7%A2%D7%A5.jpg" height="320px"/>
  <br/><span>CC BY-SA 4.0, Rena Xiaxiu, K-Pop Culture Magazine</span>
</p>

# Concept


# Build

This project is split into two separate applications. On one side the `sesamefs` application to run the actual filesystem and on the other side the client. The former can only be applied to modern Linux and FreeBSD systems, whereas the client can be compiled on any platform from Windows to macOS.

Requirements:
* Go >= 1.15
* Git

```
git clone https://github.com/outofbits/sesameFS.git
```

## Client



```
cd sesame-client && go mod vendor && go build
```


## Filesystem

```
cd sesamefs && go mod vendor && go build
```