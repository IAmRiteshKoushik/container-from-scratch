package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// -- GOAL:
// docker run image-name <cmd> <params>
// -- EXPECTED RESULT:
// go run main.go run <params>

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("bad command")
	}
}

func run() {

	// In-order to run the run() func followed by the child() func, we are /
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// Setup defaults
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,

		// NOTE: While trying this out in linux (Arch), I ran into the problem where
		//  spawning a new namespace was not permitted by the system unless being
		//  run in root mode. When I tried adding `syscall.CLONE_NEWUSER` then the
		//  program ran but if you try to run sys-admin level commands then they
		//  fail (eg: hostname container >> would not change the hostname to contianer)
		//  --
		//  In order to combat that, I am using SYS_CAP_ADMIN, running go build
		//  and providing the binary file with the required priviledges for proceeding
		//  This provides system admin level priviledges to the binary file generated
		//  which means, that by using `sudo hostname container` you can change the
		//  hostname
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("Exited gracefully")

	// WARNING: You cannot set the name of the namespace after cmd.Run() and
	//  neither can you do it before cmd.Run(), it has to be done when cmd.Run()
	//  is being executed. More details in the README.
	// syscall.Sethostname([]byte("container"))
}

func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	// Here, when we are inside the child we need not setup a separate process
	// but we do need to setup the hostname beforehand. This time it should
	// already be in the new namespace
	syscall.Sethostname([]byte("container"))

	// NOTE:
	// When a child is created, it needs to point to a separate filesystem
	// that needs to be created. So there should be some kind of
	// /root -> /container-volume-root
	// By default it points to the global system root directory, but that needs
	// to be changed. So, we create a separate filesystem manually and change the
	// directory to that
	syscall.Chroot("/home/rk/ubuntu-fs")
	// After you have done chroot, it is undefined where the root directory is
	// so you need to manually do that part
	syscall.Chdir("/")
	// In-order to run the `ps` command, we need to mount the proc directory
	// because we are in a chroot environment ?
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
	syscall.Unmount("/proc", 0)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
