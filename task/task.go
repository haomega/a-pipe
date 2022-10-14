package task

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
)

type Task struct {
	Name    string
	Request Request
}

type Request struct {
	// POST /api/xxx
	Method  string
	Path    string
	Headers []string
	Body    string
}

func LoadTaskConfig(taskName string) Task {
	fmt.Println("load", taskName)
	taskKey := "tasks." + taskName
	taskName = viper.GetString(taskKey)

	api := viper.GetString(taskKey + ".request.api")
	headers := viper.GetStringSlice(taskKey + ".request.headers")
	body := viper.GetString(taskKey + ".request.body")
	return Task{
		Name: taskName,
		Request: Request{
			Method:  strings.Split(api, " ")[0],
			Path:    strings.Split(api, " ")[1],
			Headers: headers,
			Body:    body,
		},
	}

}

func HttpRequest(task *Task) {
	fmt.Println("Run task : ", task.Name)
	req := task.Request
	url := GetBaseUrl() + req.Path

	client := &http.Client{}

	reader := strings.NewReader(req.Body)
	request, _ := http.NewRequest(req.Method, url, reader)
	// set headers
	for _, value := range req.Headers {
		split := strings.Split(value, ":")
		request.Header.Set(split[0], split[1])
	}
	// do request
	resp, _ := client.Do(request)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode, string(body))
}
