package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func resolveHost(host string) (string, error) {
	ips, err := net.LookupHost(host) // ips provide ipv6 and ipv4 address ips[0]-> ipv6 / ips[1] = ipv4
	if err != nil {
		fmt.Printf("Error resolving the hostname: %v\n", err)
		return "", err
	}

	// fmt.Println("IP addresses for", host, ":")
	// for _, ip := range ips {
	// 	fmt.Println(ip)
	// }

	return ips[1], nil
}

func sendICMPWithTTL(destination string, ttl int) (*icmp.Message, net.Addr, error) {
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0") // raw ICMP socket for IPv4
	if err != nil {
		log.Printf("listen err, %s", err)
		return nil, nil, err
	}

	defer c.Close()

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}

	wb, err := wm.Marshal(nil)
	if err != nil {
		return nil, nil, err
	}

	c.IPv4PacketConn().SetTTL(ttl)
	c.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true)

	dest := &net.IPAddr{
		IP: net.ParseIP(destination),
	}

	if _, err := c.WriteTo(wb, dest); err != nil {
		log.Printf("WriteTo err, %s\n", err)
		return nil, nil, err
	}

	replyBuffer := make([]byte, 1500)
	n, peer, err := c.ReadFrom(replyBuffer)
	if err != nil {
		return nil, nil, err
	}

	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), replyBuffer[:n])
	if err != nil {
		log.Fatal(err)
	}

	return rm, peer, nil
}

func main() {
	inputPtr := flag.String("hostname", "", "Hostname to trace route to")
	flag.Parse()
	inputHost := *inputPtr

	host, err := resolveHost(inputHost)
	if err != nil {
		log.Fatal("Error resolving host", err)
	}

	fmt.Printf("traceroute to %s(%s), 64 hops max, 32 bytes message\n", inputHost, host)
	ttl := 1
	for {
		startTime := time.Now()
		icmpResponse, peer, err := sendICMPWithTTL(host, ttl)
		if err != nil {
			panic(err)
		}
		endTime := time.Now()
		latency := endTime.Sub(startTime).Milliseconds()
		fmt.Printf("%d. %s (%d ms)\n", ttl, peer, latency)

		if icmpResponse.Type == ipv4.ICMPTypeTimeExceeded {
			ttl += 1
		}

		if icmpResponse.Type == ipv4.ICMPTypeEchoReply || ttl >= 64 {
			break
		}
	}
}
