## minikube ISO image

This includes the configuration for an alternative bootable ISO image meant to be used in conjection with minikube.

It includes:
- systemd as the init system
- rkt
- docker

## Configurations

The following configurations have been tested:

* OSX
  * Virtualbox
  * VMware Fusion

* Linux
  * Virtualbox
  * KVM

The following configurations are known to have issues currently:

* OSX
  * xhyve (https://github.com/zchee/docker-machine-driver-xhyve/issues/135)

## Hacking

### Requirements

* Linux
```
sudo apt-get install build-essential gnupg2 p7zip-full git wget cpio python \
	unzip bc gcc-multilib automake libtool locales
```

Either import your private key or generate a sign-only key using `gpg2 --gen-key`.
Also be sure to have an UTF-8 locale set up in order to build the ISO.

### Build instructions

```
$ git clone https://github.com/kubernetes/minikube
$ cd minikube
$ make buildroot-image
$ make out/minikube.iso
```

The build will occurs inside a docker container, if you want to do this
baremetal, replace `make out/minikube.iso` with `IN_DOCKER=1 make out/minikube.iso`.
The bootable ISO image will be available in `out/minikube.iso`.

### Testing local minikube-iso changes

```
$ ./out/minikube start \
    --container-runtime=rkt \
    --network-plugin=cni \
    --iso-url=file:///$GOPATH/github.com/hasura/minikube/out/buildroot/output/images/rootfs.iso9660
```

When testing a custom image and something goes bad, you will have to delete the
VM and start it again. The start command does not delete the VM if the creation
fails, and if you don't do it, you will end up starting the same incorrect
VM over and over.

```
$ ./out/minikube start --iso-url=file://...
# a wild error appears...

$ ./out/minikube delete
Deleting local Kubernetes cluster...
Machine deleted.

$ ./out/minikube start --iso-url=file://...
```


### Buildroot configuration

To change the buildroot configuration, execute:

```
$ cd out/buildroot
$ make menuconfig
$ make
```

To change the kernel configuration, execute:

```
$ cd out/buildroot
$ make linux-menuconfig
$ make
```

The last commands copies changes made to the kernel configuration to the minikube-iso defconfig.

### Saving buildroot/kernel configuration changes

To save any buildroot configuration changes made with `make menuconfig`, execute:

```
$ cd out/buildroot
$ make savedefconfig
```

The changes will be reflected in the `minikube-iso/configs/minikube_defconfig` file.

```
$ git stat
## master
 M deploy/iso/minikube-iso/configs/minikube_defconfig
```

To save any kernel configuration changes made with `make linux-menuconfig`, execute:

```
$ cd out/buildroot
$ make linux-savedefconfig
$ cp output/build/linux-4.7.2/defconfig \
    ../../deploy/iso/minikube-iso/board/coreos/minikube/linux-4.7_defconfig
```

The changes will be reflected in the `deploy/iso/minikube-iso/configs/minikube_defconfig` file.
