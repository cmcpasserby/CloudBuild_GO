package main

import (
	"fmt"
	"github.com/cmcpasserby/CloudBuild_GO/cmd/cloudbuild/cli"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	if val, ok := cli.Commands[os.Args[1]]; ok {
		err := val.Action(os.Args[2:]...)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("%q is not a valid command\n", os.Args[1])
		fmt.Println()
		printHelp()
	}
}

func printHelp() {
	fmt.Println(
		`Tool for working with Unity Cloud Build

usage:
  CloudBuild-Go <command>

commands are:`)

	for _, key := range cli.CommandOrder {
		fmt.Printf("  %-12s%s\n", cli.Commands[key].Name, cli.Commands[key].HelpText)
	}
}
