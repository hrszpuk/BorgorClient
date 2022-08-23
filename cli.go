package main

import (
	"flag"
	"fmt"
	"strings"

	"bytespace.network/rpsclient/pm"
	"bytespace.network/rpsclient/print"
)

/* cli.go handles flags and command line arguments for the project
 * Everything is documented below, this was moved away from its own package
 * formally cli/cli.go because it makes more sense for the flags to be processed
 * inside the main package - you gain no benefit from having cli as its own package.
 */

// Defaults
// These are the default values of all flags
// The values are set using the flag library, but they are also commented below
var helpFlag bool    // false -h
var showVersion bool // false -v
var options []string

// Constants that are used throughout code
// Should be updated when necessary
const executableName string = "rps"                          // in case we change it later
const discordInvite string = "https://discord.gg/kk9MsnABdF" // infinite link
const currentVersion string = "1.0"

// Init initializes and processes (parses) compiler flags
func Init() {
	flag.BoolVar(&helpFlag, "h", false, "Shows this help message")
	flag.BoolVar(&showVersion, "v", false, "Shows current package manager version")
	flag.Parse()

	options = flag.Args()
}

// ProcessFlags goes through each flag and decides how they have an effect on the output of the compiler
func ProcessFlags() {
	// Show version has higher priority than help menu
	if showVersion {
		Version()
		return

	} else if helpFlag || len(options) == 0 {
		// If they use "-h" or only enter the executable name "rps"
		// Show the help menu because they're obviously insane.
		Help()
		return
	}

	// otherwise get the operation
	op := options[0]
	if op == "get" {
		pm.Get(options[1:])
	} else if op == "update" {
		pm.Update(options[1:])
	} else if op == "remove" {
		pm.Remove(options[1:])
	} else {
		print.PrintCF(print.Red, "Unknown operation '%s'", op)
	}
}

// Help shows help message (pretty standard nothing special)
func Help() {
	header := "ReCT Package System Cli v" + currentVersion
	lines := strings.Repeat("-", len(header))

	fmt.Println(lines)
	fmt.Println(header)
	fmt.Println(lines)

	fmt.Print("\nUsage: ")
	print.PrintC(print.Green, "rps <operation> <package>\n")
	fmt.Println(print.Format("<operation> can be &g'get'&!, &y'update'&! or &r'remove'&!", print.Reset))
	fmt.Println("<package> can be the name of a package")
	fmt.Println("\n[Options]")

	helpSegments := []HelpSegment{
		{"Help", executableName + " -h", "disabled (default)", "Shows this help message!"},
		{"Version", executableName + " -v", "disabled (default)", "Shows this applications version information"},
	}

	p0, p1, p2, p3 := findPaddings(helpSegments)

	for _, segment := range helpSegments {
		segment.Print(p0, p1, p2, p3)
	}

	fmt.Println("")
	print.WriteC(print.Gray, "Still having troubles? Get help on the offical Discord server: ")
	print.WriteCF(print.DarkBlue, "%s!\n", discordInvite) // Moved so link is now blue
}

// Version Shows the current compiler version
func Version() {
	header := "ReCT Package System Cli v" + currentVersion
	lines := strings.Repeat("-", len(header))

	fmt.Println(lines)
	fmt.Println(header)
	fmt.Println(lines)

	fmt.Print("RPS version: ")
	print.PrintC(print.Blue, currentVersion)
	print.WriteC(print.Gray, "\nStill having troubles? Get help on the offical Discord server: ")
	print.WriteCF(print.DarkBlue, "%s!\n", discordInvite) // Moved so link is now blue
}

type HelpSegment struct {
	Command      string
	Example      string
	DefaultValue string
	Explanation  string
}

func (seg *HelpSegment) Print(p0 int, p1 int, p2 int, p3 int) {
	print.WriteCF(print.Cyan, "%-*s", p0, seg.Command)
	print.WriteC(print.DarkGray, ":")
	print.WriteCF(print.Blue, " %-*s", p1, seg.Example)
	print.WriteC(print.DarkGray, ":")
	print.WriteCF(print.Yellow, " %-*s", p2, seg.DefaultValue)
	print.WriteC(print.DarkGray, ":")
	print.WriteCF(print.Green, " %-*s", p3, seg.Explanation)
	fmt.Println("")
}

func findPaddings(segments []HelpSegment) (int, int, int, int) {
	p0 := 0
	p1 := 0
	p2 := 0
	p3 := 0

	for _, segment := range segments {
		if len(segment.Command) > p0 {
			p0 = len(segment.Command)
		}
		if len(segment.Example) > p1 {
			p1 = len(segment.Example)
		}
		if len(segment.DefaultValue) > p2 {
			p2 = len(segment.DefaultValue)
		}
		if len(segment.Explanation) > p3 {
			p3 = len(segment.Explanation)
		}
	}

	return p0 + 1, p1 + 1, p2 + 1, p3 + 1
}
