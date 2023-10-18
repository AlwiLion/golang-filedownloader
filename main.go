package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	urlArry := readUrlArray()

	fmt.Println("ARRAY :", urlArry)
	fmt.Println("leng ", len(urlArry)-1)
	for i := 0; i < len(urlArry); i++ {
		wg.Add(1)
		go downloadFile2(urlArry[i])
		wg.Wait()
	}

}
func readUrlArray() []string {
	fmt.Printf("How many file you want download: ")
	var size int
	fmt.Scanln(&size)
	var arr = make([]string, size)
	for i := 0; i < size; i++ {
		fmt.Printf("Enter your Url %d: ", i+1)
		var url string
		inputReader := bufio.NewReader(os.Stdin)
		url, _ = inputReader.ReadString('\n')

		url = strings.TrimSuffix(url, "\n")
		arr[i] = url
		//fmt.Println("Length Of Array : ", len(arr))
		//arr = append(arr, url)

	}
	//fmt.Println("Your array is: ", arr)
	//fmt.Println("ARRAY LENGTH: ", len(arr))
	for i := 0; i < len(arr); i++ {
		fmt.Println("Object :", i, " Value :", arr[i], "-")
	}
	return arr
}
func downloadFile2(fullURLFile string) error {

	defer wg.Done()
	// Build fileName from fullPath
	fmt.Println("---", fullURLFile)
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	//fmt.Println("Path :", path)
	segments := strings.Split(path, "/")
	//fmt.Println("Segments :", segments)
	fileName := segments[len(segments)-1]
	//fmt.Println("filename :", fileName)

	// Create blank file
	if fileExists(fileName) {
		fmt.Println("File Exist")
		tmp := strings.Split(fileName, ".")
		fileName = tmp[0] + "-" + time.Now().Format("20060102150405") + "-." + tmp[1]
	}

	fmt.Println("filename ", fileName)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("err in creating file")
		log.Fatal(err)

	}
	// Put content on file
	resp, err := http.Get(fullURLFile)
	if err != nil {
		log.Fatal(err)
	} else if resp.StatusCode != 200 {
		fmt.Println("Invalid Status Code: ")
		log.Fatal(resp.StatusCode)
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error in writting")
		log.Fatal(err)
	}

	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d", fileName, size)
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
