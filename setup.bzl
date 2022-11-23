load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")
load("@build_bazel_integration_testing//tools:repositories.bzl", "bazel_binaries")
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("//:version.bzl", "SUPPORTED_BAZEL_VERSIONS")
load("//gazelle:deps.bzl", _go_repositories = "gazelle_deps")
load("@bazel_toolchains//rules:rbe_repo.bzl", "rbe_autoconfig")

def rules_scala_internal_setup():

    # Depend on the Bazel binaries for running bazel-in-bazel tests
    bazel_binaries(versions = SUPPORTED_BAZEL_VERSIONS)

    # Bazel 5.3.0 has bzlmod bugs so we use 6.0 prerelease for the bzlmod example.
    # SUPPORTED_BAZEL_VERSIONS doesn't currently support multiple versions. For now,
    # we only want to run the bzlmod example with a separate version.
    bazel_binaries(versions = [
        "6.0.0rc1",
    ])

    bazel_skylib_workspace()

    # Creates toolchain configuration for remote execution with BuildKite CI
    # for rbe_ubuntu1604
    rbe_autoconfig(
        name = "buildkite_config",
    )

    # gazelle:repository_macro gazelle/deps.bzl%gazelle_deps
    _go_repositories()

    go_rules_dependencies()

    go_register_toolchains(version = "1.19.2")

    gazelle_dependencies()