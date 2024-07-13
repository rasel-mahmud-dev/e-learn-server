package main

func a() {

}

//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"net/http"
//	"sync"
//	"time"
//)
//
//var count int
//var mu sync.Mutex
//
//func main() {
//	var wg sync.WaitGroup
//
//	// Number of Send calls
//	numCalls := 1000
//	wg.Add(numCalls)
//	start := time.Now()
//
//	for i := 0; i < numCalls; i++ {
//		go func() {
//			defer wg.Done()
//			Send()
//		}()
//	}
//
//	wg.Wait()
//	duration := time.Since(start)
//	fmt.Println(count)
//	fmt.Println("Total Time:", duration)
//}
//
//func Send() {
//	mu.Lock()
//	defer mu.Unlock()
//
//	url := "https://api.learnwithsumit.com/v1/auth/login"
//	payload := map[string]string{
//		"email":       "admin2@test.com",
//		"password":    "173451D@",
//		"fingerprint": "7120975110",
//	}
//
//	jsonData, err := json.Marshal(payload)
//	if err != nil {
//		fmt.Println("Error marshaling JSON:", err)
//		return
//	}
//
//	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
//	if err != nil {
//		fmt.Println("Error creating request:", err)
//		return
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Accept", "application/json")
//	req.Header.Set("User-Agent", "TestAgent/1.0")
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println("Error sending request:", err)
//		return
//	}
//	defer resp.Body.Close()
//	count++
//
//	//fmt.Println(err, resp.StatusCode, resp.Status)
//}
