// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package parse

import (
	"fmt"
	"maps"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/promql/parser"
)

// As of Prometheus v0.311 the parser no longer accepts a per-instance function
// map (parser.WithFunctions was removed); custom functions can only live in the
// global parser.Functions map. Register the x-functions there unconditionally so
// parsing is deterministic regardless of how many engines exist. Per-engine
// gating (the EnableXFunctions option) is enforced above the parser, in the
// engine, after parsing.
func init() {
	maps.Copy(parser.Functions, XFunctions)
}

var XFunctions = map[string]*parser.Function{
	"xdelta": {
		Name:       "xdelta",
		ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix},
		ReturnType: parser.ValueTypeVector,
	},
	"xincrease": {
		Name:       "xincrease",
		ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix},
		ReturnType: parser.ValueTypeVector,
	},
	"xrate": {
		Name:       "xrate",
		ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix},
		ReturnType: parser.ValueTypeVector,
	},
}

// IsExtFunction is a convenience function to determine whether extended range calculations are required.
func IsExtFunction(functionName string) bool {
	_, ok := XFunctions[functionName]
	return ok
}

func UnknownFunctionError(name string) error {
	msg := fmt.Sprintf("unknown function: %s", name)
	if _, ok := parser.Functions[name]; ok {
		return errors.Wrap(ErrNotImplemented, msg)
	}

	return errors.Wrap(ErrNotSupportedExpr, msg)
}
