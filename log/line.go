package log

import (
	"strings"

	"github.com/fatih/color"
)

// Index of a log line split by blanks
const (
	dateIndex = iota
	timeIndex
	pidIndex
	tidIndex
	priorityIndex
	tagIndex
	messageIndex
)

// Line is a struct type representing a line returned from logcat
type Line struct {
	Time     string
	PID      string
	Tid      string
	Priority string
	Tag      string
	Message  string
}

// PriorityColor returns a color.Color representing a priority color (foreground/background)
func (l Line) PriorityColor() *color.Color {
	switch l.Priority {
	case "V":
		return color.New(color.FgWhite, color.BgBlack).Add(color.Bold)
	case "D":
		return color.New(color.FgBlack, color.BgBlue)
	case "I":
		return color.New(color.FgBlack, color.BgGreen)
	case "W":
		return color.New(color.FgBlack, color.BgYellow)
	case "E":
		return color.New(color.FgBlack, color.BgRed)
	case "F":
		return color.New(color.FgBlack, color.BgRed)
	default:
		return color.New()
	}
}

// NewLine parses a logcat line string and returns a Line struct representation of the fields
func NewLine(raw string) *Line {
	var time, pid, tid, tag, priority, message string

	// split line in a sequence of words
	words := strings.Fields(raw)

	// the index of the first word on the tag is always fixed
	if len(words) > tagIndex {
		// loop to build the tag until we find a ':' char, considering malformed tags or tags with spaces
		for i := tagIndex; i < len(words); i++ {
			// check for case where the ':' character is split from the tag
			if words[i] == ":" {
				message = strings.Join(words[i+1:], " ")
				break
			}

			tag += words[i]

			// check for case where ':' is appended at the end of the tag
			if len(tag) > 0 && tag[len(tag)-1] == ':' {
				tag = tag[:len(tag)-1]
				message = strings.Join(words[i+1:], " ")
				break
			}

			// add a space in-between words
			tag += " "
		}
		tag = strings.TrimSpace(tag)
	}

	if timeIndex < len(words) {
		time = words[timeIndex]
	}

	if pidIndex < len(words) {
		pid = words[pidIndex]
	}

	if tidIndex < len(words) {
		tid = words[tidIndex]
	}

	if priorityIndex < len(words) {
		priority = words[priorityIndex]
	}

	return &Line{Time: time, PID: pid, Tid: tid, Priority: priority, Tag: tag, Message: message}
}
