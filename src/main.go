package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
	"golang.org/x/sync/semaphore"
)

// PortScanner it is a struct 
type PortScanner struct{
	ip string
	lock *semaphore.Weighted
}

// Ulimit runs ulimit in bash shell and return the max open files allowed
func Ulimit() int64 {
	out, err := exec.Command("/bin/sh","-c","ulimit -n").Output()    
	if err != nil {
        panic(err)
	}  
	// fmt.Println(string(out))  
	s := strings.TrimSpace(string(out))
    i, err := strconv.ParseInt(s, 10, 64)
    
    if err != nil {
        panic(err)
	}    
	return i
}

// ScanPort function scans the port of a ip
func ScanPort(ip string, port int, timeout time.Duration){
	target := fmt.Sprintf("%s:%d",ip,port)
	conn,err := net.DialTimeout("tcp",target,timeout)

	if err != nil {
		if strings.Contains(err.Error(),"too many open files") {
			time.Sleep(timeout)
			ScanPort(ip,port,timeout)
		} else{
			fmt.Println(port,"closed")
			// if port is closed then we don't need to close the connection, hence return
			return
		}
	}
	conn.Close()
	fmt.Println(port,"open")
}

// Start is a method for struct PortScanner which will scan the port in range [f:l] along with a timeout value
func (ps *PortScanner) Start(f int, l int, timeout time.Duration){
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port:= f ; port <= l; port++ {
		wg.Add(1)
		fmt.Println()
		ps.lock.Acquire(context.TODO(),1)

		go func(port int){
			defer ps.lock.Release(1)
			defer wg.Done()
			ScanPort(ps.ip,port,timeout)
		}(port)
	}
}

func main(){
	fmt.Println(Ulimit());
	ps := &PortScanner{
		ip: "127.0.0.1",
		lock: semaphore.NewWeighted(Ulimit()),
	}
	ps.Start(1,65535,500*time.Millisecond)
}