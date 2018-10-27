package models

type WebhookCallback struct {
	PushData struct {
		Tag string
	} `json:"push_data"`

	Repository struct {
		Name     string
		RepoName string `json:"repo_name"`
	} `json:"repository"`
}
