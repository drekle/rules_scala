package scala

import (
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// GenerateRules extracts build metadata from source files in a directory.
// GenerateRules is called in each directory where an update is requested
// in depth-first post-order.
func (scala *Scala) GenerateRules(args language.GenerateArgs) language.GenerateResult {

	var result language.GenerateResult
	result.Gen = make([]*rule.Rule, 0)
	return result
}
