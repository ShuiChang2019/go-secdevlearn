package main

import (
	"SCANNER-MultiThread-Type2/scanner"
	"SCANNER-MultiThread-Type2/utils"
	"fmt"
	"os"
)

// 多线程-使用chan通道的TCP全连接端口扫描器

func main() {
	if len(os.Args) == 3 {
		ipList := os.Args[1]
		portList := os.Args[2]
		ips, err := utils.GetIpList(ipList)
		ports, err := utils.GetPorts(portList)
		_ = err

		task, _ := scanner.GenerateTask(ips, ports)
		scanner.AssigningTasks(task)
		scanner.PrintResult()

	} else {
		fmt.Printf("HELP: %v iplist port\n", os.Args[0])
	}
}
