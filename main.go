package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/pflag"
	"golang.org/x/text/unicode/runenames"
)

var version string

func main() {

	displayHelp := pflag.BoolP("help", "h", false, "display help")
	displayVersion := pflag.BoolP("version", "V", false, "display version")
	pflag.Parse()

	if *displayVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// override the default usage display
	if *displayHelp {
		displayUsage()
		os.Exit(0)
	}

	var r *bufio.Reader

	if len(pflag.Args()) > 0 {
		inputStr := strings.Join(pflag.Args(), "")
		r = bufio.NewReader(strings.NewReader(inputStr))
	} else { // arguments provided attempt to read from stdin

		// no file provided so attempt to read piped data from stdin
		info, err := os.Stdin.Stat()
		if err != nil {
			log.Fatal(err)
		}
		// ensure stdin is a pipe, bail if not
		if info.Mode()&os.ModeNamedPipe == 0 {
			log.Fatalf("Stdin isn't a pipe")
		}
		r = bufio.NewReader(os.Stdin)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"UTF-32", "Glyph", "Name"})
	table.SetBorder(false)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetColumnSeparator("")

	for {
		curRune, _, err := r.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		data := []string{fmt.Sprintf("%U", curRune), string(curRune), runenames.Name(curRune)}
		table.Append(data)
	}

	table.Render()

}

// print custom usage instead of the default provided by pflag
func displayUsage() {

	fmt.Printf("Usage: uniname [<flags>] <unicode glyphs>\n\n")
	fmt.Printf("Example: uniname 😺\n\n")
	fmt.Printf("Optional flags:\n\n")
	pflag.PrintDefaults()
}
