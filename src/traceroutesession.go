package main

import "time"
import "fmt"

/* This struct represents a traceroute session */
type TraceRouteSession struct {
	MaxTTL                  int
	Timeout                 time.Duration
	LocalAddr               string
	RemoteAddr              string
	Ipv6                    bool
	latency_estimation      time.Duration
	max_latency_between_hop time.Duration
	latency_between_hop     time.Duration
	hop_one                 int
	hop_two                 int
	hop_num                 int
}

/* initialize a traceroute session */
func NewTraceroute(localAddr string, remoteAddr string, maxTTL int, timeout int, isIpv6 bool) *TraceRouteSession {

	if localAddr == "" { /* in case the local address is empty we set a default one */
		if isIpv6 { /* ipv6 case */
			localAddr = "::"
		} else { /* ipv4 case */
			localAddr = "0.0.0.0"
		}
	}

	return &TraceRouteSession{
		MaxTTL:     maxTTL,
		Timeout:    time.Duration(timeout) * time.Second,
		LocalAddr:  localAddr,
		RemoteAddr: remoteAddr,
		Ipv6:       isIpv6,
		hop_num:    0,
		hop_one:    -1,
		hop_two:    -1,
	}
}

/* this method takes care of storing latency metrics */
func (tr *TraceRouteSession) UpdateLatencyMetrics(current_latency_estimation time.Duration, current_hop int) {
	if tr.latency_estimation == 0 { /* first hop --> just update the value */
		tr.latency_estimation = current_latency_estimation
	} else { /* update the variable holding the maximum latency between two consecutive hops */

		/* latency between current hop and the previous one */
		/* avoid negative time duration */
		tr.latency_between_hop = tr.latency_estimation - current_latency_estimation

		if tr.latency_between_hop < 0 { /* normalize negative values */
			tr.latency_between_hop *= -1
		}

		tr.latency_estimation = current_latency_estimation /* update latency time for the next hop */

		if tr.max_latency_between_hop == 0 { /* set for the first time */
			tr.max_latency_between_hop = tr.latency_between_hop
			tr.hop_one = current_hop - 1
			tr.hop_two = current_hop
		} else {
			if tr.latency_between_hop > tr.max_latency_between_hop {
				tr.max_latency_between_hop = tr.latency_between_hop
				tr.hop_one = current_hop - 1
				tr.hop_two = current_hop
			}
		}
	}
}

func (tr *TraceRouteSession) printLatencyReport() {
	if tr.max_latency_between_hop == 0 {
		fmt.Printf("No information about maximum response time between consecutive hops!\n")
	} else {
		if tr.hop_num > 1 {
			fmt.Printf("Maximum response time between consecutive hops: %s <--> Hops: [ %d - %d ]\n", tr.max_latency_between_hop, tr.hop_one, tr.hop_two)
		} else {
			fmt.Printf("No information about maximum response time between consecutive hops!\n")
		}
	}
}

/* getter method for hops num */
func (tr *TraceRouteSession) getHopsNum() int {
	return tr.hop_num
}


func (tr *TraceRouteSession) runTrace() {
	if tr.Ipv6 {
		fmt.Printf("Starting IPv6 tracing (TTL: %d, Timeout: %s)\n", tr.MaxTTL, tr.Timeout)
		tr.traceRouteIpv6()
	} else {
		fmt.Printf("Starting IPv4 tracing (TTL: %d, Timeout: %s)\n", tr.MaxTTL, tr.Timeout)
		tr.traceRouteIpv4()
	}
}
