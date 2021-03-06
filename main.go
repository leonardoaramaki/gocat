package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/leonardoaramaki/gocat/adb"
	"github.com/leonardoaramaki/gocat/log"
	"github.com/leonardoaramaki/gocat/ps"
)

// arrayFlags define a type to receive flags that can occur multiple time
type arrayFlags []string

func (i *arrayFlags) String() string {
	return "arrayFlags"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// pidOf executes 'ps' on a shell spawned at given device
func pidOf(packageName string, devices string) string {
	pid := ""
	adb.Run(devices, func(output string) {
		proc := ps.NewProcess(output)
		if proc.Package == packageName {
			pid = proc.ID
		}
	}, "ps | grep "+packageName)
	return pid
}

func main() {
	// application id or package name flag
	var packageName string
	flag.StringVar(&packageName, "p", "", "")
	flag.StringVar(&packageName, "packageName", "", "")

	// log message only with no metadata
	var raw bool
	flag.BoolVar(&raw, "r", false, "")
	flag.BoolVar(&raw, "raw", false, "")

	// ignored tags
	var ignore arrayFlags
	flag.Var(&ignore, "i", "")
	flag.Var(&ignore, "ignore", "")

	// filter by filters
	var filters arrayFlags
	flag.Var(&filters, "t", "")
	flag.Var(&filters, "tag", "")

	// copy & paste friendly
	var cp bool
	flag.BoolVar(&cp, "cp", false, "")

	// run over the emulator
	var emu bool
	flag.BoolVar(&emu, "e", false, "")
	flag.BoolVar(&emu, "emu", false, "")

	// run over usb
	var usb bool
	flag.BoolVar(&usb, "d", false, "")
	flag.BoolVar(&usb, "dev", false, "")

	// current focused app
	var current bool
	flag.BoolVar(&current, "current", false, "")

	flag.Usage = func() {
		h := "Filter logcat by package name\n\n"

		h += "Usage:\n"
		h += "	gocat -p [packageName]\n\n"

		h += "Options:\n"
		h += "	-p, --package <packageName>  Set package name to filter by\n"
		h += "	-r, --raw                    Show messages only, no metadata\n"
		h += "	-t, --tag <tag>              Filter messages with specified tag\n"
		h += "	-i, --ignore <tag>           Ignore messages with specified tag\n"
		h += "	-e --emu                     Use first emulator (adb -e)\n"
		h += "	-d --dev                     Use first device (adb -d)\n"
		h += "	-cp                          Copy & paste friendly format\n"
		h += "	--current                    Filter by current application\n\n"

		h += "Examples:\n"
		h += "	gocat -p com.example.app -i EGL_emulation -i System\n"
		h += "	gocat -p com.example.app -cp\n"

		fmt.Fprintf(os.Stderr, h)
	}

	flag.Parse()

	// Select device to connect
	d := make([]string, 0)
	if usb {
		d = append(d, "-d")
	}

	if emu {
		d = append(d, "-e")
	}

	devices := strings.Join(d, " ")

	if current {
		adb.Run(devices, func(o string) {
			if len(o) > 0 {
				splits := strings.Split(o, "/")
				if len(splits) > 0 {
					s := splits[0]
					splits = strings.Split(s, " ")
					if len(splits) > 0 {
						packageName = splits[len(splits)-1]
					}
				}
			}
		}, `dumpsys activity activities | grep mFocusedWindow`)
	}

	// id of the process running given application if any
	pid := pidOf(packageName, devices)

	// default tag width
	var tagWidth int = 23
	if raw {
		tagWidth = 0
	}

	// indentation format
	indent := fmt.Sprintf("%%%ss ", strconv.Itoa(tagWidth))

	// tag for the last line printed on the log
	lastTag := ""

	// map acommodating the tag on its entry keys to skip from logging
	ignoredTags := make(map[string]string)
	for _, tag := range ignore {
		ignoredTags[tag] = tag
	}

	// map acommodating the tag(s) to filter
	filteredTags := make(map[string]string)
	for _, tag := range filters {
		filteredTags[tag] = tag
	}

	knownTags := make(map[string]*color.Color)
	lastUsed := []*color.Color{
		color.New(color.FgRed),
		color.New(color.FgGreen),
		color.New(color.FgYellow),
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
		color.New(color.FgCyan),
	}

	brand := adb.GetProp(devices, "ro.product.manufacturer")
	sdk := adb.GetProp(devices, "ro.build.version.sdk")
	serialno := adb.GetProp(devices, "ro.serialno")
	abi := adb.GetProp(devices, "ro.product.cpu.abi")

	adb.Run(devices, func(output string) {
		var tag, prio, message string

		line := log.NewLine(output)
		tag = line.Tag
		message = strings.TrimSpace(line.Message) + "\n"

		// leaving app
		if tag == "ActivityManager" && strings.HasPrefix(line.Message, "Killing "+pid) {
			pid = pidOf(devices, packageName)
			for pid == "" {
				pid = pidOf(packageName, devices)
			}
		}

		// app got killed
		if tag == "Process" && line.Message == "Sending signal. PID: "+pid+" SIG: 9" {
			pid = pidOf(devices, packageName)
			for pid == "" {
				pid = pidOf(packageName, devices)
			}
		}

		if line.PID != pid {
			return
		}

		if !raw {
			if len(tag) > 0 && tag[len(tag)-1] == ':' {
				tag = tag[:len(tag)-1]
			}

			if len(tag) > tagWidth {
				tag = tag[:tagWidth]
			}

			prio = " " + line.Priority + " "
		}

		if pid != "" {
			skip := false
			if len(filteredTags) > 0 {
				_, filter := filteredTags[tag]
				if !filter {
					skip = true
				}
			} else {
				_, ignore := ignoredTags[tag]
				if ignore {
					skip = true
				}
			}

			if !skip {
				// truncate the tag if is the same as last one
				if lastTag == tag {
					tag = ""
				} else {
					lastTag = tag
				}

				if _, ok := knownTags[tag]; !ok {
					knownTags[tag] = lastUsed[0]
				}
				tagColor := knownTags[tag]
				lastUsed = lastUsed[1:]
				lastUsed = append(lastUsed, tagColor)

				// if copy and paste friendly
				if cp {
					if tag != "" {
						fmt.Printf("\n???  ")
						tagColor.Printf("%s ", tag)
						c := line.PriorityColor()
						c.Printf("%s", prio)
						fmt.Printf(" [%s][%s][%s][%s][%s] ", packageName, brand, sdk, serialno, abi)
						fmt.Printf("\n\n")
					}
					fmt.Printf("%s", message)
				} else {
					tagColor.Printf(indent, tag)
					c := line.PriorityColor()
					c.Printf("%s", prio)
					fmt.Printf(" %s", message)
				}
			}
		} else {
			// if no pid is set print all lines
			fmt.Printf(message)
		}
	}, "logcat")
}
