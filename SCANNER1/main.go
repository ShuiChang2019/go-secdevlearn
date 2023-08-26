package main

import (
	"SCANNER1/scanner"
	"SCANNER1/utils"
	"fmt"
	"os"
)

// 单线程TCP全连接端口扫描器

func main() {
	if len(os.Args) == 3 {
		ipList := os.Args[1]
		portList := os.Args[2]
		ips, err := utils.GetIpList(ipList)
		ports, err := utils.GetPorts(portList)
		_ = err

		for _, ip := range ips {
			for _, port := range ports {
				_, err := scanner.Connect(ip.String(), port)
				if err != nil {
					continue
				}
				fmt.Printf("ip:%v, port:%v is open\n", ip, port)
			}
		}
	} else {
		fmt.Printf("%v iplist port\n", os.Args[0])
	}
}
