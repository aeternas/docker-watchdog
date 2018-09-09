package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	flagPort = flag.String("port", "8082", "Port to listen on")
)

type WebhookCallback struct {
	PushData struct {
		Tag string `json:"tag"`
	} `json:"push_data"`

	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var webhook WebhookCallback
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error readall")
		}
		if err := json.Unmarshal(body, &webhook); err != nil {
			fmt.Printf("There was an error encoding the json. err = %s", err)
		}
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Fprint(w, "POST done")
		fmt.Println(webhook.PushData.Tag)
		fmt.Println(webhook.Repository.Name)
		//		cmdStr := "docker run -v ~/exp/a.out:/a.out ubuntu:14.04 /a.out -m 10m"
		//		out, _ := exec.Command("/bin/sh", "-c", cmdStr).Output()
		//		fmt.Printf("%s", out)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	flag.Parse()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/post", PostHandler)

	log.Printf("listening on port %s", *flagPort)
	log.Fatal(http.ListenAndServe(":"+*flagPort, mux))
}
