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

## Threat Models

* `cardano-node` is comprised, e.g. remote code execution 0-day. As long as the vulnerability doesn't make it possible to read out memory occupied by the running `cardano-node`, then sesameFS can protect against this threat.
* host system of the running `cardano-node` is comprised (root or non-root access), and the attacker is eaves dropping on the sesameFS filesystem. The filesystem itself is logging each access to the three files, which is why the operator can detect conspicuous behaviour as long as the attacker is not manipulating/erasing them. Hemce, sesameFS cannot defend against attacker with root-access that are unnoticed by the operator.
* employees directly at the cloud provider access the VPS/bare metal server and aim to read out the operational key details. In-memory storage of vault and a short life span of the throw-away key (meaning it is consumed quickly by the `cardano-node`) makes it harder for those attackers, but there is no chance of defending against them.

## Conclusion

This project doesn't aim to be a solid protection of the certification details, but an improvement over the common practise of keeping them on the host filesystem of VPS/bare metal server over which the operator has no physical/administrative control.

A big disadvantage of sesameFS is that for each restarting of the block producer, the operator has to send a key such that `cardano-node` can read the operational key details, which is why scripts checking the health and restarting the block producer automatically would not work properly. Those scripts have to assume that those key details are on the host filesystem at any time.

**Is it worth it?** I don't know :> the harm of those key details being exposed is minor compared to secret keys for pledge and node certificate.

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