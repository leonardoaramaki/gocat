# gocat

A colored logcat, based on [pidcat](https://github.com/JakeWharton/pidcat), written in Go.

## Install


You can can download a prebuilt [binary](https://github.com/leonardoaramaki/gocat/releases) for Linux, Mac or FreeBSD, or if 
you are a go developer:

```
▶ go install github.com/leonardoaramaki/gocat@latest

```

## Basic Usage

Filter log by package:

```
▶ gocat -p com.example.app.android
```

## Detailed Usage

The help message show all the possible flags, mostly the ones available on pidcat as well as a new one `-cp`. 
The `-cp` flag displays the log in a way it's easier to copy from the terminal.

```
▶ gocat --help

Filter logcat by package name

Usage:
	gocat -p [packageName]

Options:
	-p, --package <packageName>  Set package name to filter by
	-r, --raw                    Show messages only, no metadata
	-t, --tag <tag>              Filter messages with specified tag
	-i, --ignore <tag>           Ignore messages with specified tag
	-e --emu                     Use first emulator (adb -e)
	-d --dev                     Use first device (adb -d)
	-cp                          Copy & paste friendly format
    	--current                    Filter by current application

Examples:
	gocat -p com.example.app -i EGL_emulation -i System
    gocat -p com.example.app -cp
```

