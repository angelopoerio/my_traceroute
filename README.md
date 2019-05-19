# Introduction
This package contains a full working implementation of traceroute written in golang. It supports both IPv4 and IPv6.
It works in the following way:

* an ICMP echo request gets sent to the destination host, the TTL (or HopLimit for the IPv6 case) is controlled and increased at each step
* the program then waits for an ICMP answer, the following conditions can happen:
	* ICMP packet of type "Time Exceeded" is received -> a new hop has been discovered
	* ICMP packet of type "Echo Reply" is received.   -> final destination is reached, the traceroute process can end
	* no answer (a network timeout expires), this is probably an hop not behaving in standard mode (maybe a misconfigured router) 

At each step the latency is measured and showed in the output of the program, the latency measured at each step is then used to report
the response time between consecutive hops at the end of the process.

# Requirements

* golang installation (everything was developed under golang 1.9)
* MacOSx or Linux with root access
	* root access is required because the program needs to create raw sockets and this operation is privileged
	* on Linux is possible to use the extended capabilities to only grant the required permission without the need 
	  to run the program as root. The required capability is **CAP_NET_RAW**
* Permissive firewall rules for the ICMP traffic 
* IPv6 enabled connection in case you want to run an ipv6 traceroute

# Structure of the project
The project is a standard golang program, it consists of several modules:

* main.go 				- program entry point
* traceroutesession.go 	- definition of the struct representing a traceroute session & associated utility functions
* ipv4.go 				- all the ipv4 related networking code
* ipv6.go 				- all the ipv6 related networking code
* utils.go 				- utility functions definition (just a dns resolution helper)
* main_test.go 			- test suite for unit testing

No external dependencies are required, the project only depends on standard modules that can be found in the standard installation 
of the Golang compiler / runtime.
The current directory is a git repo as well, so it is possible to explore the commits history issuing the command "git log -v"

# How to build
It is enough to **cd src** in the directory of the project and issue the command **go build -o traceroute**. This will produce a statically linked executable named **traceroute** that it's ready to be used. No external libraries are required by the tool to run

# How to run the unit tests
It is enough to **cd src** in the directory of the project and issue the command **go test**.
IMPORTANT: under the hood the tests will trigger a real traceroute session so the right permissions are required to be successful

# How to run the tool
The tool can be run as the following example (ipv4):

```bash
MacBook-Pro-4:my_traceroute an.poerio$ sudo ./traceroute  -remote www.google.com 
Password:
Starting IPv4 tracing (TTL: 20, Timeout: 2s)
www.google.com resolved to 216.58.207.132, using this ipv4 address for tracing
1          192.168.178.1             2.243577ms
2         62.245.142.129            17.785228ms
3         62.245.142.128            17.177867ms
4          93.104.240.55            14.379065ms
5         108.170.247.97            16.113465ms
6         209.85.252.209            15.851289ms
7         216.58.207.132            15.924071ms
Maximum response time between consecutive hops: 15.541651ms <--> Hops: [ 1 - 2 ]
MacBook-Pro-4:my_traceroute an.poerio$ 
```

IPv6 example:

```bash
MacBook-Pro-4:my_traceroute an.poerio$ sudo ./traceroute  -remote www.google.com -ipv6
Starting IPv6 tracing (TTL: 20, Timeout: 2s)
www.google.com resolved to 2a00:1450:4016:80c::2004, using this ipv6 address for tracing
1      2001:a61:354f:5101:a96:d7ff:fe17:84f6                                 4.175788ms
2                         2001:a60::89:0:1:2                                37.953162ms
3                          2001:a60:0:106::2                                22.339593ms
4                        2001:4860:0:110c::1                                23.609172ms
5                        2001:4860:0:1::265f                                14.874908ms
6                   2a00:1450:4016:80c::2004                                16.744302ms
Maximum response time between consecutive hops: 33.777374ms <--> Hops: [ 1 - 2 ]
MacBook-Pro-4:my_traceroute an.poerio$ 
```

Please use the **-help** command line flag for more options:

```bash
MacBook-Pro-4:my_traceroute an.poerio$ ./traceroute -help
Usage of ./traceroute:
  -ipv6
    	enable ipv6 tracing (default: false)
  -local string
    	local bind address
  -remote string
    	Remote host (can be an ip or hostname)
  -timeout int
    	timeout to wait for an ICMP answer (default 2)
  -ttl int
    	ttl to use (default: 20) (default 20)
MacBook-Pro-4:my_traceroute an.poerio$ 
```

# Author
Angelo Poerio - <angelo.poerio@gmail.com>
