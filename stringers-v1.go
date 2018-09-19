package main

import "fmt"

type IPAddr [4]byte

func (ip IPAddr) String() string {
	var tmp string

	for _, val := range ip {
		tmp += "." + fmt.Sprintf("%d", val)
	}

	return tmp[1:]
}

func main() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}
}
