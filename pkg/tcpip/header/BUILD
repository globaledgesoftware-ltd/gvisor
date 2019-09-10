load("@io_bazel_rules_go//go:def.bzl", "go_test")

package(licenses = ["notice"])

load("//tools/go_stateify:defs.bzl", "go_library")

go_library(
    name = "header",
    srcs = [
        "arp.go",
        "checksum.go",
        "eth.go",
        "gue.go",
        "icmpv4.go",
        "icmpv6.go",
        "interfaces.go",
        "ipv4.go",
        "ipv6.go",
        "ipv6_fragment.go",
        "ndp_na.go",
        "ndp_ns.go",
        "ndp_opts.go",
        "tcp.go",
        "udp.go",
    ],
    importpath = "gvisor.dev/gvisor/pkg/tcpip/header",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/tcpip",
        "//pkg/tcpip/buffer",
        "//pkg/tcpip/seqnum",
        "@com_github_google_btree//:go_default_library",
    ],
)

go_test(
    name = "header_x_test",
    size = "small",
    srcs = [
        "ipversion_test.go",
        "tcp_test.go",
    ],
    deps = [":header"],
)

go_test(
    name = "header_test",
    size = "small",
    srcs = ["ndp_test.go"],
    embed = [":header"],
    deps = ["//pkg/tcpip"],
)
