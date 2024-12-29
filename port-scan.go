package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "port-scan",
		Short: "port-scan is a tool for scanning ports",
		Long:  "port-scan is a tool for scanning ports",
	}

	cmd.AddCommand(portScanCmd())

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}

func portScanCmd() *cobra.Command {
	var ip string
	var startPort int
	var endPort int

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan ports on a given IP address",
		Long:  "Scan ports on a given IP address",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Scanning IP: %s from port %d to port %d\n", ip, startPort, endPort)
			openPorts := scanPorts(ip, startPort, endPort, 500*time.Millisecond, 1000)
			if len(openPorts) == 0 {
				fmt.Println("没有开放的端口")
			} else {
				fmt.Println("端口\t状态\t运行的服务")
				fmt.Println("----------------------------")
				for _, portInfo := range openPorts {
					status := "closed"
					if portInfo.Port != 0 {
						status = "open"
					}
					fmt.Printf("%d\t%s\t%s\n", portInfo.Port, status, portInfo.Service)
				}
			}
		},
	}

	cmd.Flags().StringVarP(&ip, "ip", "i", " ", "IP address to scan")
	cmd.Flags().IntVarP(&startPort, "start-port", "s", 1, "Start port")
	cmd.Flags().IntVarP(&endPort, "end-port", "e", 65535, "End port")
	cmd.MarkFlagRequired("ip")

	return cmd
}

type PortInfo struct {
	Port    int
	Service string
}

func scanPorts(ip string, startPort, endPort int, timeout time.Duration, maxWorkers int) []PortInfo {
	ports := make(chan int, maxWorkers)
	results := make(chan PortInfo, maxWorkers)

	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range ports {
				scanPort(ip, port, timeout, results)
			}
		}()
	}

	go func() {
		for port := startPort; port <= endPort; port++ {
			ports <- port
		}
		close(ports)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var openPorts []PortInfo
	for portInfo := range results {
		if portInfo.Port != 0 {
			openPorts = append(openPorts, portInfo)
		}
	}
	sort.Slice(openPorts, func(i, j int) bool { return openPorts[i].Port < openPorts[j].Port })
	return openPorts
}

func scanPort(ip string, port int, timeout time.Duration, results chan<- PortInfo) {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err == nil {
		defer conn.Close()
		service := identifyService(conn)
		results <- PortInfo{Port: port, Service: service}
	} else {
		results <- PortInfo{Port: 0}
	}
}

func identifyService(conn net.Conn) string {
	address := conn.RemoteAddr().String()
	port := strings.Split(address, ":")[1]

	cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%s", port))
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, port) {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				return fields[0]
			}
		}
	}

	return "unknown"
}
