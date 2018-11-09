package main

import (
	"flag"
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
  CloudBuild-Go <command> [flags]
  Global Flags: --apiKey, --orgId

commands are:`)

	for _, key := range cli.CommandOrder {
		cmd := cli.Commands[key]
		fmt.Printf("  %-12s%s   flags: [", cmd.Name, cmd.HelpText)

		cmd.Flags.VisitAll(func(flag *flag.Flag) {
			if flag.Name != "apiKey" && flag.Name != "orgId" {
				fmt.Printf("--%s, ", flag.Name)
			}
		})
		fmt.Printf("\033[2D]")
		fmt.Println()
	}
}
