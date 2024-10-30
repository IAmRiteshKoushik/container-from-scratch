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
	default:
		panic("bad command")
	}
}

func run() {
	// You'll see the first argument that was the result of the go run main.go
	// fmt.Printf("%v\n", os.Args[0])

	fmt.Printf("Running %v\n", os.Args[2:])
	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	// Setup defaults
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,

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

func must(err error) {
	if err != nil {
		panic(err)
	}
}
