1. Resolving Hostname to IP Address
   
        func resolveHost(host string) (string, error) {
          ips, err := net.LookupHost(host)
          if err != nil {
            fmt.Printf("Error resolving the hostname: %v\n", err)
            return "", err
          }
          return ips[1], nil
        }
   
- This function takes a hostname (e.g., "google.com") and converts it into an IP address.
- net.LookupHost(host) retrieves a list of IPs.
- If an error occurs, it prints an error message and returns "" (empty string).
- Otherwise, it returns the second resolved IP. (first resolved ip is ipv6 and second resolved IP is ipv4)
2. Sending ICMP Packet with TTL
  
        func sendICMPWithTTL(destination string, ttl int) (*icmp.Message, net.Addr, error) {
          c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
          if err != nil {
            log.Printf("listen err, %s", err)
            return nil, nil, err
          }
          defer c.Close()
  
- `icmp.ListenPacket("ip4:icmp", "0.0.0.0")` creates a raw ICMP socket for IPv4.
- If an error occurs, it logs and returns.
- `defer c.Close()` ensures the connection is closed after execution.
3. Constructing an ICMP Echo Request

      	wm := icmp.Message{
      		Type: ipv4.ICMPTypeEcho, Code: 0,
      		Body: &icmp.Echo{
      			ID: os.Getpid() & 0xffff, Seq: 1,
      			Data: []byte("HELLO-R-U-THERE"),
      		},
      	}
  
- This creates an ICMP Echo Request (ping message).
- ID is set to the process ID (ensures uniqueness).
- `Seq: 1` is a sequence number (incremental in real traceroute).
- Data is a test message.
4. Marshaling the ICMP Message

      	wb, err := wm.Marshal(nil)
      	if err != nil {
      		return nil, nil, err
        }

- Converts the ICMP message into bytes for sending.
- If marshaling fails, the function returns an error.
5. Setting TTL and Sending the Packet

      	c.IPv4PacketConn().SetTTL(ttl)
      	c.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true)
- SetTTL(ttl): Sets the Time-To-Live for the packet.
- SetControlMessage(ipv4.FlagTTL, true): Ensures TTL changes are considered.

      	dest := &net.IPAddr{
      		IP: net.ParseIP(destination),
      	}
      
      	if _, err := c.WriteTo(wb, dest); err != nil {
      		log.Printf("WriteTo err, %s\n", err)
      		return nil, nil, err
      	}
- `net.ParseIP(destination)`: Converts the IP string into an IP address object.
- `c.WriteTo(wb, dest)`: Sends the ICMP packet to the destination.
6. Reading the Response

      	replyBuffer := make([]byte, 1500)
      	n, peer, err := c.ReadFrom(replyBuffer)
      	if err != nil {
      		return nil, nil, err
      	}
- A buffer (replyBuffer) of 1500 bytes is allocated to store the response.
- `ReadFrom` waits for a reply and retrieves:
- `n`: Number of bytes received.
- `peer`: Address of the responding host.

        	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), replyBuffer[:n])
        	if err != nil {
        		log.Fatal(err)
        	}
- The received bytes are parsed into an icmp.Message.
- If parsing fails, the program exits with an error.
7. Main Function

        func main() {
        	inputPtr := flag.String("hostname", "", "Hostname to trace route to")
        	flag.Parse()
        	inputHost := *inputPtr
- Reads the hostname argument from the command line.

        	host, err := resolveHost(inputHost)
        	if err != nil {
        		log.Fatal("Error resolving host", err)
        	}
- Converts the hostname into an IP address.

        	fmt.Printf("traceroute to %s, 64 hops max, 32 bytes message\n", inputHost)
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
- Initializes ttl = 1.
- Calls sendICMPWithTTL(), measuring the latency.
- Prints the hop number, responding host, and time taken.
8. Handling Different ICMP Responses

      		if icmpResponse.Type == ipv4.ICMPTypeTimeExceeded {
      			ttl += 1
      		}
- If a router returns Time Exceeded, increase ttl and continue.

      		if icmpResponse.Type == ipv4.ICMPTypeEchoReply || ttl >= 64 {
      			break
      		}
      	}
      }
- If the final destination responds (EchoReply), stop.
- If TTL reaches 64 (maximum hops), stop.
