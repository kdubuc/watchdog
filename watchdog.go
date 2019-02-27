package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"os/exec"
	"strings"
	"bufio"
	"bytes"
	"io/ioutil"
)

func main() {

	// Parameters list
	port := flag.String("p", "8080", "Port to listen on")
	function := flag.String("f", "", "Function to invoke")
	flag.Parse()

	// Function flag required
	if *function == "" {
		log.Fatal("Function is required")
	}

	// Assign listener port (default : 8080)
	if *port == "" {
		*port = "8080"
	}

	// Assign binding address (:[PORT])
	addr := ":" + *port

	// Log start of listening
	log.Println("Listening HTTP requests on " + addr)

	// Handle HTTP requests to invoke function
	http.HandleFunc("/", func (writer http.ResponseWriter, request *http.Request) {

		// Dump request into a raw string
		requestDump, err := httputil.DumpRequest(request, true)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Incoming Query", request.URL)
	
		// Invoke function with request dump pipe into stdin
		parts := strings.Split(*function, " ")
		fork := exec.Command(parts[0], parts[1:]...)
		fork.Stdin = strings.NewReader(string(requestDump))
		stdout, err := fork.Output()
		if err != nil {
			log.Fatal(err)
		}

		// Parse response from function stdout
		response, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(stdout)), request)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("OK", response.StatusCode)

		// Write response's status code
		writer.WriteHeader(response.StatusCode)

		// Write response's body
		var body, _ = ioutil.ReadAll(response.Body)
		writer.Write(body)

		// Write response's headers
		for key, values := range response.Header {
			writer.Header().Set(key, strings.Join(values, ";"))
		}
	
	})

	// Listen HTTP requests
	http.ListenAndServe(addr, nil)

}
