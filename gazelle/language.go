package scala

import (
	"github.com/bazelbuild/bazel-gazelle/language"
)

// Python satisfies the language.Language interface. It is the Gazelle extension
// for Python rules.
type Scala struct {
	Configurer
	Resolver
}

// NewLanguage initializes a new Python that satisfies the language.Language
// interface. This is the entrypoint for the extension initialization.
func NewLanguage() language.Language {
	return &Scala{}
}
