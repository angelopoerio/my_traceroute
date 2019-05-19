package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
	"os"
	"time"
)

/* run a traceroute for ipv4 */
func (tr *TraceRouteSession) traceRouteIpv4() {
	var dst net.IPAddr

	dnsResolve(tr.RemoteAddr, false, &dst)

	/* istantiate a listening socket for ICMP answers */
	icmp_sock, err := net.ListenPacket("ip4:icmp", tr.LocalAddr)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not set a listening ICMP socket: %s\n", err)
		os.Exit(1)
	}

	/* clean up the socket when the function exits */
	defer icmp_sock.Close()

	/* istantiate a raw IPv4 socket because we need to manipulate the TTL field in the header */
	ipv4_sock := ipv4.NewPacketConn(icmp_sock)
	defer ipv4_sock.Close()

	if err := ipv4_sock.SetControlMessage(ipv4.FlagTTL|ipv4.FlagDst|ipv4.FlagInterface|ipv4.FlagSrc, true); err != nil {
		fmt.Fprintf(os.Stderr, "Could not set options on the ipv4 socket: %s\n", err)
		os.Exit(1)
	}

	/* istantiate an ICMP echo request used in the traceroute process, payload is just an empty string because we do not need it */
	icmp_echo := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0, Body: &icmp.Echo{ID: rand.Int(), Data: []byte("")},
	}

	/* usual MTU size */
	read_buf := make([]byte, 1500)

	/* main logic is contained below */
	for i := 1; i < tr.MaxTTL; i++ {
		icmp_echo.Body.(*icmp.Echo).Seq = i

		write_buffer, err := icmp_echo.Marshal(nil)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not serialize the ICMP echo request: %s\n", err)
			os.Exit(1)
		}

		/* set TTL for the current Hop */
		if err := ipv4_sock.SetTTL(i); err != nil {
			fmt.Fprintf(os.Stderr, "Could not set the TTL field: %s\n", err)
			os.Exit(1)
		}

		time_now := time.Now()

		if _, err := ipv4_sock.WriteTo(write_buffer, nil, &dst); err != nil {
			fmt.Fprintf(os.Stderr, "Could not send the ICMP packet: %s\n", err)
			os.Exit(1)
		}

		/* set timeout to avoid to have a forever blocking call */
		if err := ipv4_sock.SetReadDeadline(time.Now().Add(tr.Timeout)); err != nil {
			fmt.Fprintf(os.Stderr, "Could not set the read timeout on the ipv4 socket: %s\n", err)
			os.Exit(1)
		}

		read_bytes, _, hop_node, err := ipv4_sock.ReadFrom(read_buf)

		/* network error or timeout, just skip and go to the next hop */
		if err != nil {
			fmt.Printf("%d %20s\n", i, "*")
		} else { /* got an answer */
			icmp_answer, err := icmp.ParseMessage(1, read_buf[:read_bytes])

			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not parse the ICMP packet from: %s\n", hop_node.String())
				os.Exit(1)
			}

			current_latency_estimation := time.Since(time_now)

			tr.UpdateLatencyMetrics(current_latency_estimation, i)

			if icmp_answer.Type == ipv4.ICMPTypeTimeExceeded {
				fmt.Printf("%d   %20s   %20s\n", i, hop_node.String(), current_latency_estimation) /* new hop */
				tr.hop_num++
			} else if icmp_answer.Type == ipv4.ICMPTypeEchoReply {
				fmt.Printf("%d   %20s   %20s\n", i, hop_node.String(), current_latency_estimation) /* final host */
				break
			} else {
				fmt.Printf("%d %20s\n", i, "*") /* other ICMP packets we are not interested in */
			}

		}
	}
}