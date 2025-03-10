package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

func getURLString(port string) string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Error getting hostname:", err)
		return fmt.Sprintf("using port %s with password %s", port, password)
	}

	addresses, err := net.LookupIP(hostname)
	if err != nil {
		log.Println("Error looking up IP addresses:", err)
		return fmt.Sprintf("using port %s with password %s", port, password)
	}

	for _, addr := range addresses {
		if ipv4 := addr.To4(); ipv4 != nil {
			//fmt.Println("IPv4:", ipv4.String())
			if disableSSL {
				return fmt.Sprintf("http://%s:%s@%s:%s", username, password, ipv4.String(), port)
			} else {
				return fmt.Sprintf("https:/%s:%s@%s:%s", username, password, ipv4.String(), port)
			}
		}
	}
	log.Fatalln("no host ip address found")
	return ""
}

func generatePassword() {
	const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 12)
	for i := range b {
		b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	password = string(b)
}
