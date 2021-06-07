package tinylog

import (
	"fmt"
	"io"
	"os"

	"github.com/andriiyaremenko/tinylog/formatters"
)

// Returns all Destinations that were assigned to factory.
func AllDestinations(factory LoggerFactory) []Destination {
	return factory.Destinations()
}

// Destination constructor function.
func DestinationFunc(out io.Writer, formatter formatters.LogFormatter, level int) Destination {
	return func() *destination { return &destination{out: out, formatter: formatter, level: level} }
}

// Destination based on os.Stderr as out and formatters.Default() as formatter.
func DefaultDestination() *destination {
	return &destination{out: os.Stderr, formatter: formatters.Default(), level: Info}
}

type destination struct {
	level     int
	out       io.Writer
	formatter formatters.LogFormatter
}

func (d *destination) ID() string {
	return fmt.Sprintf("%T->%T<%p>", d.formatter, d.out, d.out)
}
