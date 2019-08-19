package main

import (
	"flag"
	"os"
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
)

func main() {
	opt := godog.Options{Output: colors.Colored(os.Stdout)}
	godog.BindFlags("godog.", flag.CommandLine, &opt)
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		err := FeatureContext(s)
		if err != nil {
		    fmt.Fprintf(os.Stderr, "Initialise Error\n%#v", err)
			os.Exit(1)
		}
	}, opt)

	os.Exit(status)
}
