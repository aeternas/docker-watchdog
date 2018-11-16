package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	. "github.com/aeternas/docker-watchdog/models"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

var (
	flagPort = flag.String("port", "8083", "Port to listen on")
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var webhook WebhookCallback
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("error readall")
		}
		if err := json.Unmarshal(body, &webhook); err != nil {
			log.Printf("There was an error encoding the json. err = %s", err)
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		fmt.Fprint(w, "POST done")
		log.Println(webhook.PushData.Tag)
		log.Println(webhook.Repository.RepoName)
		log.Println(webhook.Repository.Name)
		imageName := fmt.Sprintf("%s:%s", webhook.Repository.RepoName, webhook.PushData.Tag)
		log.Println(imageName)
		cmdStr := fmt.Sprintf("./docker_update.sh %s %s %s", webhook.Repository.RepoName, webhook.PushData.Tag, webhook.Repository.Name)
		log.Println(cmdStr)
		out, err := exec.Command("/bin/bash", "-c", cmdStr).Output()
		if err != nil {
			log.Println("error!")
			text := fmt.Sprintf("{\"attachments\":[{\"fallback\":\"Failed to deploy application\",\"color\":\"#ff0000\",\"author_name\":\"Docker-Watchdog\",\"title\":\"Deploy result\",\"text\":\"Failed to deploy application: %s:%s\"}]}", webhook.Repository.RepoName, webhook.PushData.Tag)
			deploymentResult(text)
		} else {
			text := fmt.Sprintf("{\"attachments\":[{\"fallback\":\"Application has been successfully deployed\",\"color\":\"#36a64f\",\"author_name\":\"Docker-Watchdog\",\"title\":\"Deploy result\",\"text\":\"Application has been successfully deployed: %s:%s\"}]}", webhook.Repository.RepoName, webhook.PushData.Tag)
			deploymentResult(text)
		}
		log.Println(out)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func deploymentResult(s string) {
	url := "https://hooks.slack.com/services/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"
	fmt.Println("URL:>", url)

	log.Println(s)
	var jsonStr = []byte(s)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	client.Timeout = time.Second * 15
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(body)
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
