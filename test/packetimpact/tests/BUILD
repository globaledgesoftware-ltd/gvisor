load("defs.bzl", "packetimpact_go_test")

package(
    default_visibility = ["//test/packetimpact:__subpackages__"],
    licenses = ["notice"],
)

packetimpact_go_test(
    name = "syn_with_same_ip",
    srcs = ["syn_with_same_ip_test.go"],
    deps = [
        "//pkg/tcpip/header",
        "//test/packetimpact/testbench",
        "@com_github_mohae_deepcopy//:go_default_library",
        "@org_golang_x_sys//unix:go_default_library",
    ],
)

sh_binary(
    name = "test_runner",
    srcs = ["test_runner.sh"],
)
