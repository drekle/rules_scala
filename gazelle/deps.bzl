"This file managed by `bazel run //:update_go_deps`"

load("@bazel_gazelle//:deps.bzl", _go_repository = "go_repository")

def go_repository(name, **kwargs):
    if name not in native.existing_rules():
        _go_repository(name = name, **kwargs)

def gazelle_deps():
    "Fetch go dependencies"
    go_repository(
        name = "com_github_bazelbuild_buildtools",
        build_naming_convention = "go_default_library",
        importpath = "github.com/bazelbuild/buildtools",
        sum = "h1:jhiMzJ+8unnLRtV8rpbWBFE9pFNzIqgUTyZU5aA++w8=",
        version = "v0.0.0-20221004120235-7186f635531b",
    )