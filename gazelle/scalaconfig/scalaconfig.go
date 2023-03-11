package scalaconfig

import (
	"path/filepath"
	"strings"
)

// Directives
const (
	// ScalaExtensionDirective represents the directive that controls whether
	// this Scala extension is enabled or not. Sub-packages inherit this value.
	// Can be either "enabled" or "disabled". Defaults to "enabled".
	ScalaExtensionDirective = "scala_extension"
	// ScalaRootDirective represents the directive that sets a Bazel package as
	// a Scala root. This is used on monorepos with multiple Scala projects
	// that don't share the top-level of the workspace as the root.
	ScalaRootDirective = "scala_root"

	packageNameNamingConventionSubstitution = "$package_name$"
)

// Configs is an extension of map[string]*Config. It provides finding methods
// on top of the mapping.
type Configs map[string]*Config

// ParentForPackage returns the parent Config for the given Bazel package.
func (c *Configs) ParentForPackage(pkg string) *Config {
	dir := filepath.Dir(pkg)
	if dir == "." {
		dir = ""
	}
	parent := (map[string]*Config)(*c)[dir]
	return parent
}

// Config represents a config extension for a specific Bazel package.
type Config struct {
	parent                  *Config
	scalaProjectRoot        string
	libraryNamingConvention string
}

// New creates a new Config.
func New(
	repoRoot string,
	scalaProjectRoot string,
) *Config {
	return &Config{
		scalaProjectRoot:        scalaProjectRoot,
		libraryNamingConvention: packageNameNamingConventionSubstitution,
	}
}

// Parent returns the parent config.
func (c *Config) Parent() *Config {
	return c.parent
}

// NewChild creates a new child Config. It inherits desired values from the
// current Config and sets itself as the parent to the child.
func (c *Config) NewChild() *Config {
	return &Config{
		parent: c,
	}
}

// SetPythonProjectRoot sets the Python project root.
func (c *Config) SetScalaProjectRoot(scalaProjectRoot string) {
	c.scalaProjectRoot = scalaProjectRoot
}

// PythonProjectRoot returns the Python project root.
func (c *Config) ScalaProjectRoot() string {
	return c.scalaProjectRoot
}

// RenderLibraryName returns the py_library target name by performing all
// substitutions.
func (c *Config) RenderLibraryName(packageName string) string {
	return strings.ReplaceAll(c.libraryNamingConvention, packageNameNamingConventionSubstitution, packageName)
}
