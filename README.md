# Traceroute in go

### What is Traceroute?
Traceroute is a network diagnostic tool used to trace the path that data packets take from a source computer to a destination. It helps in identifying the route taken by the packets and detecting network issues such as high latency or packet loss.

#### How Traceroute Works:

- It sends ICMP Echo Request packets (or UDP/TCP packets in some implementations) to the target.
- It starts with a low TTL (Time-To-Live) value (usually 1) and increases it step by step.
- When a router receives a packet with a TTL of 1, it discards the packet and sends back an ICMP Time Exceeded message.
- The program logs the IP address of each router along the path.
- This process continues until the packet reaches the final destination, which responds with an ICMP Echo Reply.

### Run locally

    git clone git@github.com:utsab818/traceroute-in-go.git

#### Run main.go file
    sudo go run main.go -hostname=google.com

#### Why sudo?
Sending raw ICMP packets requires root (administrator) privileges on most operating systems.

Project idea: https://github.com/aerosouund/go-tracert/tree/master
