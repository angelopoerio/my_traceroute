package main

import "testing"

/*
	Test suit  
*/

func TestIPv4Tracing(test *testing.T) {
	t := NewTraceroute("0.0.0.0", "www.google.com", 10, 2, false)
	t.runTrace()
	if t.getHopsNum() == 0 {
		test.Errorf("IPv4 traceroute: Expected more than 0 hops when tracerouting www.google.com")
	}
}

func TestIPv6Tracing(test *testing.T) {
	t := NewTraceroute("::", "www.google.com", 10, 2, true)
	t.runTrace()
	if t.getHopsNum() == 0 {
		test.Errorf("IPv6 traceroute: Expected more than 0 hops when tracerouting www.google.com")
	}
}