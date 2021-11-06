package ps

import "strings"

/*
	Process is a struct type that represents a line output from a 'ps' shell command
*/
type Process struct {
	UserId  string // User Id
	ID      string // Process Id
	Ppid    string // Parent Process Id
	Package string
}

/*
	NewProcess parses a string and return a Process struct
*/
func NewProcess(raw string) *Process {
	fields := strings.Fields(raw)

	return &Process{UserId: fields[0], ID: fields[1], Ppid: fields[2], Package: fields[8]}
}
