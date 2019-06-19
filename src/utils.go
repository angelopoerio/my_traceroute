package main

import (
	"fmt"
	"net"
	"os"
)

/* DNS resolution utility function (IPv6/IPv4) */
func dnsResolve(hostname string, isIPv6 bool, dst *net.IPAddr) {
	ip_addr := net.ParseIP(hostname)
	/* check if the provided address is an ipv6 address or an hostname */
	if isIPv6 && ip_addr.To16() != nil {
		fmt.Printf("Using the provided ipv6 address %s for tracing\n", hostname)
		dst.IP = ip_addr
		/* check if the provided address is an ipv4 address or an hostname */
	} else if !isIPv6 && ip_addr.To4() != nil {
		fmt.Printf("Using the provided ipv4 address %s for tracing\n", hostname)
		dst.IP = ip_addr
	} else {
		ips, err := net.LookupIP(hostname)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not resolve %s\n", hostname)
			os.Exit(1)
		}

		for _, ip := range ips {
			/* look for the "AAAA" record containing the ipv6 address we want to resolve */
			if isIPv6 && ip.To16() != nil {
				dst.IP = ip
				fmt.Printf("%s resolved to %s, using this ipv6 address for tracing\n", hostname, ip)
				break
				/* look for the "A" record containing the ipv4 address we want to resolve */
			} else if !isIPv6 && ip.To4() != nil {
				dst.IP = ip
				fmt.Printf("%s resolved to %s, using this ipv4 address for tracing\n", hostname, ip)
				break
			}
		}

		if dst.IP == nil {
			fmt.Printf("Could not find a valid record for %s\n", hostname)
			os.Exit(1)
		}
	}
}
