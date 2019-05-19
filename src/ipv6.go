package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv6"
	"math/rand"
	"net"
	"os"
	"time"
)

/* run a traceroute for ipv6 */
func (tr *TraceRouteSession) traceRouteIpv6()  {
	var dst net.IPAddr

	dnsResolve(tr.RemoteAddr, true, &dst)

	/* istantiate a listening socket for ICMP6 answers */
	icmp6_sock, err := net.ListenPacket("ip6:ipv6-icmp", tr.LocalAddr)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not set a listening ICMP6 socket: %s\n", err)
		os.Exit(1)
	}

	/* clean up the socket when the function exits */
	defer icmp6_sock.Close()

	/* istantiate a raw IPv6 socket because we need to manipulate the HopLimit field in the header */
	ipv6_sock := ipv6.NewPacketConn(icmp6_sock)
	defer ipv6_sock.Close()

	if err := ipv6_sock.SetControlMessage(ipv6.FlagHopLimit|ipv6.FlagDst|ipv6.FlagInterface|ipv6.FlagSrc, true); err != nil {
		fmt.Fprintf(os.Stderr, "Could not set options on the ipv6 socket: %s\n", err)
		os.Exit(1)
	}

	/* istantiate an ICMP echo request used in the traceroute process, payload is just an empty string because we do not need it */
	icmp6_echo := icmp.Message{
		Type: ipv6.ICMPTypeEchoRequest, Code: 0, Body: &icmp.Echo{ID: rand.Int(), Data: []byte("")},
	}

	/* usual MTU size */
	read_buf := make([]byte, 1500)

	/* main logic is contained below */
	for i := 1; i < tr.MaxTTL; i++ {
		icmp6_echo.Body.(*icmp.Echo).Seq = i

		write_buffer, err := icmp6_echo.Marshal(nil)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not serialize the ICMP6 echo request: %s\n", err)
			os.Exit(1)
		}

		/* set the HopLimit for the current Hop (no TTL in the header of an IPV6 packet) */
		if err := ipv6_sock.SetHopLimit(i); err != nil {
			fmt.Fprintf(os.Stderr, "Could not set the HopLimit field: %s\n", err)
			os.Exit(1)
		}

		time_now := time.Now()

		if _, err := ipv6_sock.WriteTo(write_buffer, nil, &dst); err != nil {
			fmt.Fprintf(os.Stderr, "Could not send the ICMP6 echo packet: %s\n", err)
			os.Exit(1)
		}

		/* set timeout to avoid to have a forever blocking call */
		if err := ipv6_sock.SetReadDeadline(time.Now().Add(tr.Timeout)); err != nil {
			fmt.Fprintf(os.Stderr, "Could not set the read timeout on the ipv6 socket: %s\n", err)
			os.Exit(1)
		}

		read_bytes, _, hop_node, err := ipv6_sock.ReadFrom(read_buf)

		/* network error or timeout, just skip and go to the next hop */
		if err != nil {
			fmt.Printf("%d %40s\n", i, "*")
		} else { /* got an answer */
			icmp_answer, err := icmp.ParseMessage(58, read_buf[:read_bytes])

			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not parse the ICMP6 packet from: %s\n", hop_node.String())
				os.Exit(1)
			}

			current_latency_estimation := time.Since(time_now)

			tr.UpdateLatencyMetrics(current_latency_estimation, i)

			if icmp_answer.Type == ipv6.ICMPTypeTimeExceeded {
				fmt.Printf("%d   %40s   %40s\n", i, hop_node.String(), current_latency_estimation) /* new hop */
				tr.hop_num++
			} else if icmp_answer.Type == ipv6.ICMPTypeEchoReply {
				fmt.Printf("%d   %40s   %40s\n", i, hop_node.String(), current_latency_estimation) /* final host */
				break
			} else {
				fmt.Printf("%d %40s\n", i, "*")
			}

		}
	}
}
