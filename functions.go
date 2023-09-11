package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

func enableSeDebugPrivilege() error {
	var t windows.Token
	if err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_ALL_ACCESS, &t); err != nil {
		return err
	}

	var luid windows.LUID

	if err := windows.LookupPrivilegeValue(nil, windows.StringToUTF16Ptr(seDebugPrivilege), &luid); err != nil {
		return fmt.Errorf("LookupPrivilegeValueW failed, error: %v", err)
	}

	ap := windows.Tokenprivileges{
		PrivilegeCount: 1,
	}

	ap.Privileges[0].Luid = luid
	ap.Privileges[0].Attributes = windows.SE_PRIVILEGE_ENABLED

	if err := windows.AdjustTokenPrivileges(t, false, &ap, 0, nil, nil); err != nil {
		return fmt.Errorf("AdjustTokenPrivileges failed, error: %v", err)
	}

	return nil
}

func checkIfAdmin() bool {
	f, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func elevate() error {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	if err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd); err != nil {
		return err
	}

	os.Exit(0)
	return nil
}
