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

While creating a new namespace using `syscall.CLONE_NEWNS`, I have added the 
`Unshareflags: syscall.CLONE_NEWNS` which allows me to keep the data of the 
namespace hidden. Things that are mounted inside the containerized environment
will not be visible to the host.
  
A before-after example of `mount | grep proc`
```bash
# before (observe the last line)
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
systemd-1 on /proc/sys/fs/binfmt_misc type autofs (rw,relatime,fd=39,pgrp=1,timeout=0,minproto=5,maxproto=5,direct,pipe_ino=6362)
binfmt_misc on /proc/sys/fs/binfmt_misc type binfmt_misc (rw,nosuid,nodev,noexec,relatime)
proc on /home/rk/ubuntu-fs/proc type proc (rw,relatime)

# after (the last line disappeared which is the relevant mount)
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
systemd-1 on /proc/sys/fs/binfmt_misc type autofs (rw,relatime,fd=39,pgrp=1,timeout=0,minproto=5,maxproto=5,direct,pipe_ino=6362)
binfmt_misc on /proc/sys/fs/binfmt_misc type binfmt_misc (rw,nosuid,nodev,noexec,relatime)
```

You can still look it up by putting the process to sleep and then using the pid
```bash
sudo cat /proc/<pid>/mounts
```

## Chroot (change root)
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

## Cgroups (control groups)

This tells you what you can "use". Things like the filesystem interface
- Memory
- CPU
- I/O
- Process numbers

The way control groups are accessed also differs from OS to OS. 

In the program, the we have setup a simple control group which controls the 
number of processes than can be run inside the container.

```bash
# The code limits the processes to be 20 in count, so in-order to test that 
:(){ :|: & };:
```

This particular sketchy looking command starts forking the existing process on 
a linux terminal. When that happens the max limit of 20 processes comes into play 
and blocks further forking. This is basically a fork-bomb.

> [!INFO] 
> Tampering with control groups is fairly tricky. Be careful and figure out ways 
> to remove control-groups after your job is done.

