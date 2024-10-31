# Build your own container in Go

> [!WARNING]
> This repository must be cloned and run in Linux / Unix environments only. It 
is not built for windows.

1. Namespaces
2. Chroot (changing the root)
3. Cgroups (control groups)

When you run a container using docker - `docker run --rm -it ubuntu /bin/bash`
This runs the container and then removes it once the work is done. It is by 
default executing the bash shell 

## Namespaces
Limit what a process can see. Created using `syscalls`
- Unix Timesharing System
- Process IDs
- Mounts
- Network
- User IDs
- InterProcess Comms

This plays a major role in restricting / isolating a container.

After having created a separate and isolated process

## Chroot
The chroot command can be used to setup and run programs or interactive shells 
such as Bash in an encapsulated filesystem that is prevented from interacting 
with your regular filesystem. Everything within the chroot environment is penned 
in and contained. Nothing in the chroot environment can see out past its own, 
special root directory.

Learn how to create a chroot directory here - [link](https://www.howtogeek.com/441534/how-to-use-the-chroot-command-on-linux/)  
Name this a `ubuntu-fs` and to run it in the code. The new container getting 
created will use this as the root point. This need not be the actual ubuntu 
filesystem unless and until you are running on ubuntu.

> [!NOTE]
> `chroot` creation is at times, OS dependent and corresponding OS friendly 
> guides are available on the internet. Personally, the above link did not 
> workout for me as I am using Arch but for an Ubuntu user it would work just 
> fine.

Alternative is to get a filesystem and chroot into it programatically.
```bash
# Here, I am taking a redis docker container and extracting out the entire 
# filesystem to chroot into it. Ultimately it is a linux filesystem and all the 
# basic commands work. While it is not completely isolated and the host can 
# still access it, a certain level of permission setup should do the trick.
docker export $(docker create ubuntu) -o ubuntu.tar.gz
mkdir ubuntu-fs && cd ubuntu-fs
tar --no-same-owner --no-same-permissions --owner=0 --group=0 -mxf ../ubuntu.tar.gz
```

Another alternative suggested by a friend is to use `debootstrap` program. It 
creates a debian environment inside a sub-directory and lets you access it.

## Cgroups
