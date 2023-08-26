package scanner

import (
	"SCANNER-MultiThread-Type2/config"
	"fmt"
	"net"
	"strings"
	"sync"
)

// 生成扫描任务列表

func GenerateTask(ipList []net.IP, ports []int) ([]map[string]int, int) {
	tasks := make([]map[string]int, 0)

	for _, ip := range ipList {
		for _, port := range ports {
			ipPort := map[string]int{ip.String(): port} // 绑定IP和端口
			tasks = append(tasks, ipPort)
		}
	}

	return tasks, len(tasks)
}

// 分割扫描任务，根据并发数分割成组

func AssigningTasks(tasks []map[string]int) {
	// 分为scanBatch批
	scanBatch := len(tasks) / config.ThreadNum

	for i := 0; i < scanBatch; i++ {
		curTask := tasks[config.ThreadNum*i : config.ThreadNum*(i+1)]
		RunTask(curTask)
	}

	// 剩下的没分完的内容
	if len(tasks)%config.ThreadNum > 0 {
		lastTasks := tasks[config.ThreadNum*scanBatch:]
		RunTask(lastTasks)
	}
}

func RunTask(tasks []map[string]int) {
	wg := &sync.WaitGroup{}

	// 创建一个buffer为config.threadNum * 2的channel
	taskChan := make(chan map[string]int, config.ThreadNum*2)

	// 创建config.ThreadNum个协程
	for i := 0; i < config.ThreadNum; i++ {
		go Scan(taskChan, wg)
	}

	// 生产者，不断地往taskChan channel发送数据，直接channel阻塞
	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}

	close(taskChan)
	wg.Wait()
}

func Scan(taskChan chan map[string]int, wg *sync.WaitGroup) {
	// 每个协程都从channel中读取数据后开始扫描并入库
	for task := range taskChan {
		for ip, port := range task {
			err := SaveResult(Connect(ip, port))
			_ = err
			wg.Done()
		}
	}
}

func SaveResult(ip string, port int, err error) error {
	if err != nil {
		return err
	}

	v, ok := config.Result.Load(ip)
	if ok {
		ports, ok1 := v.([]int)
		if ok1 {
			ports = append(ports, port)
			config.Result.Store(ip, ports)
		}
	} else {
		ports := make([]int, 0)
		ports = append(ports, port)
		config.Result.Store(ip, ports)
	}
	return err
}

func PrintResult() {
	config.Result.Range(func(key, value interface{}) bool {
		fmt.Printf("IP:%v\n", key)
		fmt.Printf("PORTS:%v\n", value)
		fmt.Println(strings.Repeat("-*-", 30))
		return true
	})
}
