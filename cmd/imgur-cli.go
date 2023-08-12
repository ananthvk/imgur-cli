package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/ananthvk/imgur-cli/internal"
)

var writer = flag.CommandLine.Output()

// Custom usage/help function
func Usage() {
	fmt.Fprintf(writer, "imgur-cli: A CLI to upload images to imgur anonymously\n")
	fmt.Fprintf(writer, "Get the latest version at https://github.com/ananthvk/imgur-cli\n\n")
	fmt.Fprintf(writer, "Usage:\n\n")
	flag.PrintDefaults()
	/*
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(writer, " %v\n", f.Usage)
		})
	*/
	fmt.Fprintf(writer, "The avaiable commands are:\n\n")
	fmt.Fprintf(writer, "%s\t%s\n", "upload", "Uploads the given file(s) to imgur")
	fmt.Fprintf(writer, "%s\t%s\n", "delete", "Delete the files specified by the delete hash(es)")
	fmt.Fprintf(writer, "\nExamples:\n\n")
	fmt.Fprintf(writer, "%s upload cat.png\n", os.Args[0])
	fmt.Fprintf(writer, "%s delete <Delete hash>\n", os.Args[0])
}

// Tells the user to run the command with the help flag and exit
func exitRunHelp() {
	fmt.Fprintf(writer, "Run \"%s --help\" for help\n", os.Args[0])
	os.Exit(2)
}

func Run() {
	flag.Usage = Usage
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Fprint(writer, "No arguments provided\n")
		exitRunHelp()
	}
	args := flag.Args()
	// Handle the two sub commands - upload and delete
	// Also perform basic error checking
	switch args[0] {
	case "upload":
		if len(args) == 1 {
			fmt.Fprintf(writer, "No files provided to upload to imgur\n")
			exitRunHelp()
		}
		internal.Upload(args[1:])
	case "delete":
		if len(args) == 1 {
			fmt.Fprintf(writer, "No delete hashes provided to delete from imgur\n")
			exitRunHelp()
		}
		internal.Delete(args[1:])
	default:
		fmt.Fprintf(writer, "%s : unknown command\n", args[0])
		exitRunHelp()
	}
}
