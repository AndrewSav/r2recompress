package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AndrewSav/r2recompress/internal/core"
	"github.com/AndrewSav/r2recompress/internal/model"
	"github.com/AndrewSav/r2recompress/internal/version"
)

// TODO: dump strings mode

func main() {

	fmt.Printf("r2recompress %s\n", version.BuildVersion())

	w := flag.CommandLine.Output()
	flag.Usage = func() {
		fmt.Fprintf(w, "Usage: %s [options] inputFile outputFile\n", os.Args[0])
		flag.PrintDefaults()
	}

	options := model.Options{}

	options.Decompress = flag.Bool("d", false, "decompress")
	options.Compress = flag.Bool("c", false, "compress")
	options.StringDump = flag.Bool("s", false, "dump strings from a decompressed file")
	options.SingleWarning = flag.Bool("q", false, "decompress only - only print the first warning if there are warnings")
	options.NoWarnings = flag.Bool("qq", false, "decompress only - do not print warnings")
	options.Verbose = flag.Bool("v", false, "print file debug data")

	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		fmt.Fprintf(w, "Expected 2 arguments got %d", flag.NArg())
		os.Exit(2)
	}

	commands := 0
	if *options.Decompress {
		commands++
	}
	if *options.Compress {
		commands++
	}
	if *options.StringDump {
		commands++
	}

	if commands != 1 {
		flag.Usage()
		fmt.Fprintf(w, "Exactly one of -c, -d and -s should be specified")
		os.Exit(2)
	}

	if *options.Decompress {
		err := core.Decompress(flag.Arg(0), flag.Arg(1), options)
		if err != nil {
			fmt.Printf("error decompressing: %v\n", err)
		}
	}

	if *options.Compress {
		err := core.Compress(flag.Arg(0), flag.Arg(1), options)
		if err != nil {
			fmt.Printf("error compressing: %v\n", err)
		}
	}

	if *options.StringDump {
		err := core.StringDump(flag.Arg(0), flag.Arg(1), options)
		if err != nil {
			fmt.Printf("error extracting strings: %v\n", err)
		}
	}
}
