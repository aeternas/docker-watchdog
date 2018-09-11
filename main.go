package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

var (
	flagPort = flag.String("port", "8083", "Port to listen on")
)

type WebhookCallback struct {
	PushData struct {
		Tag string `json:"tag"`
	} `json:"push_data"`

	Repository struct {
		RepoName string `json:"repo_name"`
	} `json:"repository"`
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var webhook WebhookCallback
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("error readall")
		}
		if err := json.Unmarshal(body, &webhook); err != nil {
			log.Printf("There was an error encoding the json. err = %s", err)
		}
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Fprint(w, "POST done")
		log.Println(webhook.PushData.Tag)
		log.Println(webhook.Repository.RepoName)
		imageName := fmt.Sprintf("%s:%s", webhook.Repository.RepoName, webhook.PushData.Tag)
		log.Println(imageName)
		cmdStr := fmt.Sprintf("./docker_update.sh %s %s", webhook.Repository.RepoName, webhook.PushData.Tag)
		log.Println(cmdStr)
		out, err := exec.Command("/bin/bash", "-c", cmdStr).Output()
		if err != nil {
			log.Println("error!")
		}
		log.Println(out)
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
