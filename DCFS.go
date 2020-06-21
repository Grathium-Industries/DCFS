package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// INITILIZATION
	fmt.Println("Finding Peers...")
	getServers()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		files := (r.URL.Query()).Get("file")

		// not a false positive
		if files != "" {
			log.Print("[" + files + "]\n")
			fmt.Fprintf(w, readFile("files/"+files))
		}
	})

	fs := http.FileServer(http.Dir("files/"))
	http.Handle("/files/", http.StripPrefix("/files/", fs))

	fmt.Println("\nStarted server on port :4422")
	http.ListenAndServe(":4422", nil)
}

// find p2p servers
func getServers() {
	file, err := os.Open("/files/servers.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	for _, serverListIndex := range txtlines {
		// serverListIndex = current scanning line
		serverScanning := string(getHTML(serverListIndex + "/?file=servers.txt"))

		// read the remote server list
		// and compare if the server is already known
		for _, remoteServeList := range split(serverScanning, '$') {
			// read the whole file at once
			b, err := ioutil.ReadFile("files/servers.txt")
			if err != nil {
				panic(err)
			}
			localServerList := string(b)

			// compare local list to index
			for i := 0; i < len(remoteServeList); i++ {
				if strings.Contains(localServerList, remoteServeList[i]) == false {
					writeToFile("files/servers.txt", remoteServeList[i]+"$")
				}
			}

		}
	}
}

// create remote server list array
func split(tosplit string, sep rune) []string {
	var fields []string

	last := 0
	for i, c := range tosplit {
		if c == sep {
			// Found the separator, append a slice
			fields = append(fields, string(tosplit[last:i]))
			last = i + 1
		}
	}

	// Don't forget the last field
	fields = append(fields, string(tosplit[last:]))

	return fields
}

func readFile(path string) string {
	// Open file for reading.
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return "404"
	}
	defer file.Close()

	// Read file, line by line
	var text = make([]byte, 1024)
	for {
		_, err = file.Read(text)

		// Break if finally arrived at end of file
		if err == io.EOF {
			break
		}

		// Break if error occured
		if err != nil && err != io.EOF {
			isError(err)
			break
		}
	}

	return string(text)
}

func writeToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func getHTML(server string) []byte {
	url := server
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// show the HTML code as a string %s
	return html
}

/* error checking function */
func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}
