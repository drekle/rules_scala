package scala

import (
	"flag"
	"fmt"
	"log"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/rules_scala/gazelle/scalaconfig"
)

// Configurer satisfies the config.Configurer interface. It's the
// language-specific configuration extension.
type Configurer struct{}

// RegisterFlags registers command-line flags used by the extension. This
// method is called once with the root configuration when Gazelle
// starts. RegisterFlags may set an initial values in Config.Exts. When flags
// are set, they should modify these values.
func (scala *Configurer) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {}

// CheckFlags validates the configuration after command line flags are parsed.
// This is called once with the root configuration when Gazelle starts.
// CheckFlags may set default values in flags or make implied changes.
func (scala *Configurer) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	return nil
}

// KnownDirectives returns a list of directive keys that this Configurer can
// interpret. Gazelle prints errors for directives that are not recoginized by
// any Configurer.
func (scala *Configurer) KnownDirectives() []string {
	return []string{
		scalaconfig.ScalaExtensionDirective,
		scalaconfig.ScalaRootDirective,
	}
}

// Configure modifies the configuration using directives and other information
// extracted from a build file. Configure is called in each directory.
//
// c is the configuration for the current directory. It starts out as a copy
// of the configuration for the parent directory.
//
// rel is the slash-separated relative path from the repository root to
// the current directory. It is "" for the root directory itself.
//
// f is the build file for the current directory or nil if there is no
// existing build file.
func (scala *Configurer) Configure(c *config.Config, rel string, f *rule.File) {
	// Create the root config.
	if _, exists := c.Exts[languageName]; !exists {
		rootConfig := scalaconfig.New(c.RepoRoot, "")
		c.Exts[languageName] = scalaconfig.Configs{"": rootConfig}
	}

	configs := c.Exts[languageName].(scalaconfig.Configs)

	config, exists := configs[rel]
	if !exists {
		parent := configs.ParentForPackage(rel)
		config = parent.NewChild()
		configs[rel] = config
	}

	if f == nil {
		return
	}

	// gazelleManifestFilename := "gazelle_scala.yaml"

	for _, d := range f.Directives {
		switch d.Key {
		case "exclude":
			// We record the exclude directive for coarse-grained packages
			// since we do manual tree traversal in this mode.
			// config.AddExcludedPattern(strings.TrimSpace(d.Value))
		case scalaconfig.ScalaExtensionDirective:
			switch d.Value {
			case "enabled":
				// config.SetExtensionEnabled(true)
			case "disabled":
				// config.SetExtensionEnabled(false)
			default:
				err := fmt.Errorf("invalid value for directive %q: %s: possible values are enabled/disabled",
					scalaconfig.ScalaExtensionDirective, d.Value)
				log.Fatal(err)
			}
		case scalaconfig.ScalaRootDirective:
			config.SetScalaProjectRoot(rel)
		}
	}

	// gazelleManifestPath := filepath.Join(c.RepoRoot, rel, gazelleManifestFilename)
	// gazelleManifest, err := py.loadGazelleManifest(gazelleManifestPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if gazelleManifest != nil {
	// 	config.SetGazelleManifest(gazelleManifest)
	// }
}
