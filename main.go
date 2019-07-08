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
	"os"
	"os/exec"
	"time"
)

const (
	SLACK_WEBHOOK_KEY = "SLACK_WEBHOOK"
)

var (
	flagPort = flag.String("port", "8083", "Port to listen on")
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	var webhook WebhookCallback
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error on ReadAll: ", err)
	}
	if err := json.Unmarshal(body, &webhook); err != nil {
		log.Printf("There was an error encoding the json. err = %s", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	log.Println("POST done: ", w)
	log.Println("TAG is ", webhook.PushData.Tag)
	log.Println("RepoName is ", webhook.Repository.RepoName)
	log.Println("Repository.Name is ", webhook.Repository.Name)
	imageName := fmt.Sprintf("%s:%s", webhook.Repository.RepoName, webhook.PushData.Tag)
	log.Println(imageName)
	var cmdStr string
	if webhook.Repository.Name == "swadeshness-nginx" {
		cmdStr = fmt.Sprintf("./nginx_deployment.sh %s %s %s", webhook.Repository.RepoName, webhook.PushData.Tag, webhook.Repository.Name)
	} else if webhook.Repository.Name == "swadeshness-spring" {
		cmdStr = fmt.Sprintf("./spring_deployment.sh %s %s %s", webhook.Repository.RepoName, webhook.PushData.Tag, webhook.Repository.Name)
	} else {
		cmdStr = fmt.Sprintf("./docker_update.sh %s %s %s", webhook.Repository.RepoName, webhook.PushData.Tag, webhook.Repository.Name)
	}
	log.Println(cmdStr)
	out, err := exec.Command("/bin/bash", "-c", cmdStr).Output()
	if err != nil {
		log.Println("Deploy script execution error: ", err)
		text := fmt.Sprintf("{\"attachments\":[{\"fallback\":\"Failed to deploy application\",\"color\":\"#ff0000\",\"author_name\":\"Docker-Watchdog\",\"title\":\"Deploy result\",\"text\":\"Failed to deploy application: %s:%s with error %s\"}]}", webhook.Repository.RepoName, webhook.PushData.Tag, err)
		deploymentResult(text)
	} else {
		text := fmt.Sprintf("{\"attachments\":[{\"fallback\":\"Application has been successfully deployed\",\"color\":\"#36a64f\",\"author_name\":\"Docker-Watchdog\",\"title\":\"Deploy result\",\"text\":\"Application has been successfully deployed: %s:%s\"}]}", webhook.Repository.RepoName, webhook.PushData.Tag)
		deploymentResult(text)
	}
	log.Println(out)
}

func deploymentResult(s string) {
	var slackKey string
	if value, ok := os.LookupEnv(SLACK_WEBHOOK_KEY); ok {
		slackKey = value
	}
	url := fmt.Sprintf("https://hooks.slack.com/services/%s", slackKey)
	log.Println("URL:> ", url)

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

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erro on readAll: ", err)
	}
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

	go func() {
		for {
			time.Sleep(5 * time.Second)
		}
	}()
}
