# Running Unit Tests With Containerd

> The test environment is complex for the programmer to set up.
> Fortunately, **Docker** is popular for people to learn and love as a common language for **temporarily creating test environments**.
> However, there are many alternatives for docker, such as **Containerd**, Kata Container, and gVisor.
> Most likely, there is a high probability that **Containerd** replaces Docker in **the container management ecosystem** because Containerd is **less dependent and more efficient**.
> Indeed, Containerd does not have **docker cli layer**.

## Installation

### Introduction

What is **Containerd**？

**Containerd** has been *used* for *a long time* as **a container management tool**. However, People cannot be aware of its existence because it is **a container runtime** with simplicity and does not have **docker cli layer**.

<img src="./assets/image-20220331141145598.png" alt="image-20220331141145598" style="zoom:100%;" /> 

### Install main component

> Please refer to [the installation guide](https://containerd.io/downloads/) for specific details.

install Containerd components

```bash
# libseccomp2 is a security package used to prevent attackers bypass intended access restrictions for argument-filtered system calls.
$ sudo apt-get update
$ sudo apt-get install libseccomp2

# download containerd package. current version is 1.6.2.
$ wget https://github.com/containerd/containerd/releases/download/v1.6.2/cri-containerd-cni-1.6.2-linux-amd64.tar.gz

# cri-containerd-cni includes runc and internet components.
$ tar -tf cri-containerd-cni-1.6.2-linux-amd64.tar.gz | grep runc
usr/local/bin/containerd-shim-runc-v2
usr/local/bin/containerd-shim-runc-v1
usr/local/sbin/runc # runc is a CLI tool for running containers.

# install containerd
$ sudo tar -C / -xzf cri-containerd-cni-1.6.2-linux-amd64.tar.gz

# check the existence of systemd config.
$ tar -tf cri-containerd-cni-1.6.2-linux-amd64.tar.gz | grep containerd.service
etc/systemd/system/containerd.service # systemd config.

# start Containerd Service.
$ sudo systemctl daemon-reload # reload Systemd.
$ sudo systemctl enable --now containerd.service # start Containerd service at boot.
$ sudo systemctl start containerd.service # start Containerd service.

# check the existence of ctr cli tool.
$ tar -tf cri-containerd-cni-1.6.2-linux-amd64.tar.gz | grep ctr
usr/local/bin/ctr # ctr cli tool.

# test ctr cli tool.
$ ctr container list
# display CONTAINER    IMAGE    RUNTIME
```

### Managing plugins

set up Containerd plugins

```bash
# check the existence of cni component.
$ tar -tf cri-containerd-cni-1.6.2-linux-amd64.tar.gz | grep opt/cni
# The result is displayed as follows.
opt/cni/
opt/cni/bin/
opt/cni/bin/tuning
opt/cni/bin/vrf
opt/cni/bin/loopback
opt/cni/bin/portmap
opt/cni/bin/ptp
opt/cni/bin/ipvlan
opt/cni/bin/host-device
opt/cni/bin/macvlan
opt/cni/bin/host-local
opt/cni/bin/firewall
opt/cni/bin/bandwidth
opt/cni/bin/sbr
opt/cni/bin/vlan
opt/cni/bin/static
opt/cni/bin/bridge
opt/cni/bin/dhcp

# Generate config.toml config
$ tar -tf cri-containerd-cni-1.6.2-linux-amd64.tar.gz | grep config.toml # does not contain config file.
$ containerd config default > /etc/containerd/config.toml # create default config file.
```

### Compile cni tool

```bash
# compile cnitool
$ git clone https://github.com/containernetworking/cni.git
$ cd cni
$ go mod tidy
$ cd cnitool
$ go build .

# move cnitool to bin folder
$ mv ./cnitool /usr/local/bin
```

# Subnetwork

> to force containers to be isolated with a small subnetwork of limited hosts numbers.



SoarTest is not in an isolated test environment and depends on many technic things, such as containers, databases, SQL, etc.
I should tear down those dependencies as much as possible.
Fewer dependencies make a more stable test.

I make an essential decision to integrate Containerd with UnitTest.
Containerd will take over docker sooner or later.
I choose Containerd to keep up with technology after evaluating it.
Rename from SoarTest to ContainerdTest.