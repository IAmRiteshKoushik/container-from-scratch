# Build your own container in Go

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

## Chroot

## Cgroups
