package trace

import (
	"fmt"
	"io"
)

// Tracer is the interface that describes an object capable of
// tracing events throughout code.
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

// New returns a tracer pointer
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// nilTracer sruct differs from tracer struct
// in that it doesn't take an io.Writer interface
// it doesn't need on as it's not going to write anything
type nilTracer struct{}

// nillTraer struct defines a Trace method that does nothing
func (t *nilTracer) Trace(a ...interface{}) {}

// Off creates a Tracer that will ignore calls to trace
func Off() Tracer {
	return &nilTracer{}
}
