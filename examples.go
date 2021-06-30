package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"regexp"
)

var colorMap = map[string]string{
	"BOLD":    "\u001b[1m",
	"BLACK":   "\u001b[30m",
	"RED":     "\u001b[31m",
	"GREEN":   "\u001b[32m",
	"YELLOW":  "\u001b[33m",
	"BLUE":    "\u001b[34m",
	"MAGENTA": "\u001b[35m",
	"CYAN":    "\u001b[36m",
	"WHITE":   "\u001b[37m",
	"RESET":   "\u001b[0m",
}

func printExamples() {
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println(sarColors(examplesString))
	} else {
		fmt.Println(stripColors(examplesString))
	}
}

// TODO: colors
func sarColors(a_string string) string {
	return a_string
}

func stripColors(a_string string) string {
	reColor := regexp.MustCompile("%^(BOLD|RESET|BLACK|RED|GREEN|YELLOW|BLUE|MAGENTA|CYAN|WHITE)%^")
	return reColor.ReplaceAllString(a_string, "")
}

var examplesString = `
Getting Info

	exceptions list
		List all exceptions currently in the database and not marked 
		  "deleted" (under soft-delete scheme)
		Can also list more specific categories: default is "exceptions list all"
	
	exceptions details 4
	  Print detailed information about the exception with ID "4".
		Note: "detail" and "info" can also be used here, they do 
		  the same thing.
	
Submitting a New Exception

	exceptions submit --username=ccspapp --service=grace --type=quota --detail="25TB Scratch"
	  Logs a new exception, starting today and ending in a year, for Grace, with 
			the text given in the detail option.
		Note that "type" is currently a free string, but "quota" and "queue" are expected types.

Statuses

  Exceptions are expected to go through the following statuses:
	  - undecided
		- approved or rejected (rejected stops here)
		- implemented
		- removed
	
	These are set by a command of the same name:
	  exceptions approve 4
		exceptions implemented 4
		exceptions remove 4
		exceptions reject 4
		exceptions undecide 4

	The command-line interface will try and keep you to sensible transitions, but you can add
	  "-f" to force it.
  
`
