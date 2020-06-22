/*
	This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// INITILIZATION
	fmt.Println("Finding Peers...")
	discoverServers()

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

	fmt.Println("Started server on port :4422")
	fmt.Println("============================")
	http.ListenAndServe(":4422", nil)

	// initilize and load GUI
	GUILoad()
}

// find p2p servers
func discoverServers() {
	for _, serverListIndex := range split(readFile("files/servers.txt"), '$') {
		// serverListIndex = current scanning line
		remoteServer := string(getHTML(serverListIndex + "/?file=servers.txt"))

		// read the remote server list
		// and compare if the server is already known
		for _, remoteServeList := range split(remoteServer, '$') {

			// read the whole file at once
			b, err := ioutil.ReadFile("files/servers.txt")
			if err != nil {
				panic(err)
			}
			localServerList := string(b)

			// compare local list to index
			// append if not known in local list
			if strings.Contains(localServerList, remoteServeList) == false {
				writeToFile("files/servers.txt", remoteServeList+"$")
			}

		}
	}
}

// select server to use from local server list
func selectServer() string {
	localServerCount := strings.Count(readFile("file/servers.txt"), "$")

	var localServers []string
	for _, serverListRead := range split(readFile("files/servers.txt"), '$') {
		localServers = append(localServers, serverListRead)
	}

	return localServers[randomNum(0, localServerCount)]
}

// GUILoad componants and rendering
func GUILoad() {
	for {
		fmt.Println("Finding Peers...")
		fmt.Println("Started server on port :4422")
		fmt.Println("============================")
		fmt.Println("\nGet File,")

		// get user input and load website
		webBrowserOpen(selectServer() + "/?file=" + getInput())

		// vanity function
		clearScreen()
	}
}

// open web browser
func webBrowserOpen(website string) {
}

// random number function
func randomNum(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return (rand.Intn(max-min+1) + min)
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

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

/* error checking function */
func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func clearScreen() { print("\033[H\033[2J") }
