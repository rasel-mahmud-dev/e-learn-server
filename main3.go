package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/sys/windows"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

var count int
var mu sync.Mutex

func main323() {
	// Example: Command to change IP address (Windows specific)
	//cmd := exec.Command("netsh", "interface", "ipv4", "show", "config")
	//cmd := exec.Command("netsh", "interface", "ipv4", "set", "address", "Wi-Fi 2", "static", "10.100.10.1", "255.255.255.0", "10.100.10.29")

	var wg sync.WaitGroup

	// Number of Send calls
	numCalls := 100
	wg.Add(numCalls)
	start := time.Now()

	for i := 0; i < numCalls; i++ {
		go func() {
			defer wg.Done()
			Send()
		}()
	}

	wg.Wait()
	duration := time.Since(start)
	fmt.Println(count)
	fmt.Println("Total Time:", duration)

}

func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}

func amAdmin(mu sync.Mutex) bool {
	mu.Lock()
	defer mu.Unlock()

	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("admin no")
		return false
	}

	cmd := exec.Command("netsh", "interface", "ipv4", "set", "address", "Wi-Fi 2", "static", "10.100.10.37", "255.255.255.0", "10.100.10.1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running netsh command:", err)
	}
	time.Sleep(time.Second * 30)
	fmt.Println(string(output))
	fmt.Println("admin yes")
	return true
}

func Send() {

	url := "https://api.learnwithsumit.com/v1/auth/login"
	payload := map[string]string{
		"email":       "admin2@test.com",
		"password":    "173451D@",
		"fingerprint": "7120975110",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "TestAgent/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	count++

	if resp.StatusCode == 429 {
		if !amAdmin(mu) {
			runMeElevated()
		}

	} else {
		mu.Lock()
		defer mu.Unlock()
	}
}
