package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ascheglov/belfry/pkg/belfry"
)

func main() {
	args := &belfry.RunArgs{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "belfry [-h] [-p port] [-bastion host] destination [command]\n\n")
		flag.PrintDefaults()
	}
	flag.StringVar(&args.Port, "p", "22", "remote port")
	flag.StringVar(&args.Bastion, "bastion", "", "bastion host")
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		os.Exit(255)
	}
	if flag.NArg() == 0 {
		fmt.Fprint(os.Stderr, "Error: need a destination\n\n")
		flag.Usage()
		os.Exit(255)
	}

	args.Host = flag.Arg(0)

	if flag.NArg() > 1 {
		args.Command = flag.Args()[1:]
	}

	if err := belfry.Run(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
