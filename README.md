# anyparent

This is an experimental utility for MS Windows to start a new child process under any existing, currently running parent process. The new subprocess will run under the same user as the parent process and has the same privileges.

This can be useful for all sorts of troubleshooting. For example you could open a `cmd.exe` shell that has the same privileges as a service account and interactively test things like file system access with the same privileges as the service process.

## usage

To use the tool, admin privileges (for `SeDebugPrivileges`) are required (but the started process might have less privileges, depending on its parent).

```
anyparent -p <pid> [-c] <cmdline>
```


- `-p` pid (required)
- `-c` console mode (optional, default=false): set this for console applications, including cmd.exe or powershell.exe, etc (this will start conhost, otherwise there won't be any console output)

When no `<cmdline>` is provided, a cmd.exe shell is started.

## example:

Let's start a cmd.exe shell as subprocess for this MS Paint parent process:

```batch
$ tasklist /fi "imagename eq mspaint.exe"
Image Name   PID 
============ ======
mspaint.exe  23484 
```

from an elevated shell run:
```
anyparent.exe -p 23484 -c cmd.exe
```

And we will see a `cmd.exe` shell, that is a parent process of `mspaint.exe`. It also runs under the same user and with the same privileges as `mspaint.exe` (check i.e. with `whoami /all`).

![screenshot process tree](_static/image.png)

## origin

This is a fork of the excellent [FourCoreLabs/TrustedInstallerPOC](https://github.com/FourCoreLabs/TrustedInstallerPOC.git), which is made for a specific use case (creating a specific subprocess for one specific Windows service). The original author has also written [a blog post about how it works](https://fourcore.io/blogs/no-more-access-denied-i-am-trustedinstaller).

`anyparent` makes some minor modifications to allow using this technique for different use cases. That's also why I picked a new name.

