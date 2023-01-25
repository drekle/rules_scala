package scala

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/rules_scala/gazelle/scalaconfig"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/emirpasic/gods/sets/treeset"
	godsutils "github.com/emirpasic/gods/utils"
	"github.com/google/uuid"
)

// GenerateRules extracts build metadata from source files in a directory.
// GenerateRules is called in each directory where an update is requested
// in depth-first post-order.
func (scala *Scala) GenerateRules(args language.GenerateArgs) language.GenerateResult {

	cfgs := args.Config.Exts[languageName].(scalaconfig.Configs)
	cfg := cfgs[args.Rel]

	scalaProjectRoot := cfg.ScalaProjectRoot()
	packageName := filepath.Base(args.Dir)

	scalaLibraryFilenames := treeset.NewWith(godsutils.StringComparator)

	for _, f := range args.RegularFiles {
		ext := filepath.Ext(f)
		// FIXME: This is more complicated however for early test
		// Scala files may have java dependencies
		if ext == ".scala" {
			scalaLibraryFilenames.Add(f)
		}
	}

	// Add files from subdirectories if they meet the criteria.
	for _, d := range args.Subdirs {
		// boundaryPackages represents child Bazel packages that are used as a
		// boundary to stop processing under that tree.
		boundaryPackages := make(map[string]struct{})
		err := filepath.WalkDir(
			filepath.Join(args.Dir, d),
			func(path string, entry fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				return nil
			})
		_ = boundaryPackages
		// TODO:
		// NOOP
		if err != nil {
			panic(err)
		}
	}

	parser := newScalaParser(args.Config.RepoRoot, args.Rel, cfg.IgnoresDependency)
	_ = parser
	visibility := fmt.Sprintf("//%s:__subpackages__", scalaProjectRoot)

	var result language.GenerateResult
	result.Gen = make([]*rule.Rule, 0)
	collisionErrors := singlylinkedlist.New()
	_ = collisionErrors

	if !scalaLibraryFilenames.Empty() {

		scalaLibraryTargetName := cfg.RenderLibraryName(packageName)
		// TODO: add / invoke the scala parser
		deps := treeset.NewWith(moduleComparator)

		scalaLibrary := newTargetBuilder(scalaLibraryKind, scalaLibraryTargetName, scalaProjectRoot, args.Rel).
			setUUID(uuid.Must(uuid.NewUUID()).String()).
			addVisibility(visibility).
			addSrcs(scalaLibraryFilenames).
			addModuleDependencies(deps).
			generateImportsAttribute().
			build()

		_ = scalaLibrary
		// result.Gen = append(result.Gen, scalaLibrary)
		// result.Imports = append(result.Imports, scalaLibrary.PrivateAttr(config.GazelleImportsKey))

	}

	if !collisionErrors.Empty() {
		it := collisionErrors.Iterator()
		for it.Next() {
			log.Printf("ERROR: %v\n", it.Value())
		}
		os.Exit(1)
	}

	return result
}
