package main

import (
	"flag"
	"fmt"
	"github.com/cmcpasserby/ucb/cmd/cloudbuild/cli"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	if val, ok := cli.Commands[os.Args[1]]; ok {
		flagsMap, err := cli.ParseFlags(val.Flags, os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

		err = val.Action(flagsMap)
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
  ucb <command> [flags]
  Global Flags: --apiKey, --orgId (these are best defined in the config file via 'CloudBuild-Go config')

commands are:`)

	maxNameLen := 0
	maxDescLen := 0

	for _, key := range cli.CommandOrder {
		cmd := cli.Commands[key]
		if len(cmd.Name) > maxNameLen {
			maxNameLen = len(cmd.Name)
		}

		if len(cmd.HelpText) > maxDescLen {
			maxDescLen = len(cmd.HelpText)
		}
	}
	maxNameLen += 2
	maxDescLen += 2

	for _, key := range cli.CommandOrder {
		cmd := cli.Commands[key]
		fmt.Printf("  %-*s%-*sflags: [", maxNameLen, cmd.Name, maxDescLen, cmd.HelpText)

		hasFlags := false

		cmd.Flags.VisitAll(func(flag *flag.Flag) {
			if flag.Name != "apiKey" && flag.Name != "orgId" {
				fmt.Printf("--%s, ", flag.Name)
				hasFlags = true
			}
		})

		if hasFlags {
			fmt.Printf("\033[2D]")
		} else {
			fmt.Printf("]")
		}

		fmt.Println()
	}
}
