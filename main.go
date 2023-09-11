package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

const (
	seDebugPrivilege = "SeDebugPrivilege"
)

func RunAsTrustedInstaller(ppid uint32, exe string, args []string) error {
	if !checkIfAdmin() {
		if err := elevate(); err != nil {
			return fmt.Errorf("cannot elevate Privs: %v", err)
		}
	}

	if err := enableSeDebugPrivilege(); err != nil {
		return fmt.Errorf("cannot enable %v: %v", seDebugPrivilege, err)
	}

	hand, err := windows.OpenProcess(windows.PROCESS_CREATE_PROCESS|windows.PROCESS_DUP_HANDLE|windows.PROCESS_SET_INFORMATION, true, ppid)
	if err != nil {
		return fmt.Errorf("cannot open ti process: %v", err)
	}

	cmd := exec.Command(exe, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_CONSOLE,
		ParentProcess: syscall.Handle(hand),
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("cannot start new process: %v", err)
	}

	fmt.Fprintf(os.Stdout, "Started process with PID %d using cmdline: %s %s", cmd.Process.Pid, exe, strings.Join(args[:], " "))
	return nil
}

func main() {

	flag.Usage = func() {
		exe := filepath.Base(os.Args[0])
		fmt.Printf("Usage: %s -p [parent pid (required)] [-c (optional)] [command (optional)]\n", exe)
		flag.PrintDefaults()
	}

	FlagPid := flag.Uint("p", 0, "parent pid (required)")
	FlagConsole := flag.Bool("c", false, "is console application? (= use conhost) (optional)")
	flag.Parse()
	ppid := uint32(*FlagPid)
	useConhost := *FlagConsole

	if ppid == 0 {
		fmt.Fprintln(os.Stderr, "ERROR - no PPID provided.")
		flag.Usage()
		os.Exit(1)
	}

	var exe string
	var args []string

	if len(os.Args) < 4 {
		exe = "conhost.exe"
		args = []string{}
	} else {
		if useConhost {
			exe = "conhost.exe"
		} else {
			exe = os.Args[3]
		}
		args = os.Args[4:]
	}

	if err := RunAsTrustedInstaller(ppid, exe, args); err != nil {
		panic(err)
	}
}
