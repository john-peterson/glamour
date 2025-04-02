package main

// This app render markdown and downsample colors when
// necessary per the detected color profile of the terminal.

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/glamour"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Read the data.
	var buf bytes.Buffer
	if _, readErr := buf.ReadFrom(f); readErr != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", readErr)
		os.Exit(1)
	}

	// Create a new colorprofile writer. We'll use it to detect the color
	// profile and downsample colors when necessary.
	r, w, err := os.Pipe()
	_ = w
	if err != nil {
		panic(err)
	}
	 os.Setenv("CLICOLOR_FORCE", "1")
	c := colorprofile.NewWriter(w, os.Environ())
	// c := colorprofile.NewWriter(os.Stdout, os.Environ())

	// While we're at it, let's jot down the detected color profile in the
	// markdown output while we're at it.
	fmt.Fprintf(&buf, "\n\nBy the way, this was rendererd as _%s._\n", c.Profile)

	// Okay, now let's render some markdown.
	g, err := glamour.NewTermRenderer(glamour.WithEnvironmentConfig())
	if err != nil {
		log.Fatal(err)
	}
	md, err := g.RenderBytes(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	// And finally, write it to stdout using the colorprofile writer. This will
	// ensure colors are downsampled if necessary.
	fmt.Fprintf(c, "%s\n", md)
	// panic("crash")

	// cmd := exec.Command("cat")
	cmd := exec.Command("less")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	// panic("crash")

}

