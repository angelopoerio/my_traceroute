package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	remoteAddress := flag.String("remote", "", "Remote host (can be an ip or hostname)")
	localAddress := flag.String("local", "", "local bind address")
	ttl := flag.Int("ttl", 20, "ttl to use (default: 20)")
	timeout := flag.Int("timeout", 2, "timeout to wait for an ICMP answer")
	ipv6 := flag.Bool("ipv6", false, "enable ipv6 tracing (default: false)")
	flag.Parse()

	if *remoteAddress == "" {
		fmt.Fprintf(os.Stderr, "Remote address is mandatory\n")
		os.Exit(1)
	}

	t := NewTraceroute(*localAddress, *remoteAddress, *ttl, *timeout, *ipv6)
	t.runTrace()
	t.printLatencyReport()

	os.Exit(0) /* be shell friendly, exit with a success status code */
}
