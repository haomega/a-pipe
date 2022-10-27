package task

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Task struct {
	Name    string
	Api     Api
	Headers map[string]string
	Body    Body
}

type Api struct {
	Method string
	Url    string
}

type Body struct {
	Type string
	Data string
}

func LoadTask(taskName string) (*Task, error) {
	taskKey := "tasks." + taskName
	apiKey := taskKey + ".api"
	customDomainKey := taskKey + ".domain"
	customPortsKey := taskKey + ".port"
	headersKey := taskKey + ".headers"
	bodyTypeKey := taskKey + ".body.type"
	bodyDataKey := taskKey + ".body.data"

	if viper.InConfig(taskKey) {
		api := viper.GetString(apiKey)
		headers := getHeaders(headersKey)
		bodyType := viper.GetString(bodyTypeKey)
		bodyData := viper.GetString(bodyDataKey)
		apiSplit := GetKeyValuePair(api, " ")
		if apiSplit == nil {
			return nil, errors.New("api conf illagle " + api)
		}
		method := apiSplit.Key
		url := getUrl(viper.GetString(customDomainKey), viper.GetString(customPortsKey), apiSplit.Value)

		return &Task{
			Name:    taskName,
			Api:     Api{method, url},
			Headers: headers,
			Body:    Body{bodyType, bodyData},
		}, nil
	}
	return nil, errors.New("task not found " + taskName)
}

func getUrl(customDomain string, customPort string, path string) string {
	baseConfig := LoadBaseConfig()
	if customDomain == "" {
		customDomain = baseConfig.Domain
	}
	if customPort == "" {
		customPort = baseConfig.Port
	}
	return customDomain + ":" + customPort + path
}

func (task *Task) RequestApi(parameters map[string]string) error {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	api := task.Api
	body := task.Body

	method := api.Method
	url, lessPus := expandArgs(api.Url, parameters)
	bodyData, lessPbs := expandArgs(body.Data, parameters)
	less := append(lessPus, lessPbs...)
	if len(less) > 0 {
		return errors.New("less parameters " + strings.Join(less, ","))
	}
	var request *http.Request

	switch body.Type {
	case "json":
		request = getJsonRequest(method, url, bodyData)
		break
	case "form-data":
		req, err := getFormDataRequest(method, url, bodyData)
		if err != nil {
			return err
		}
		request = req
		break
	default:
		req, _ := http.NewRequest(method, url, nil)
		request = req
	}

	// set headers
	for key := range task.Headers {
		if request.Header.Get(key) == "" {
			request.Header.Set(key, task.Headers[key])
		}
	}
	// do request
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	all, err := io.ReadAll(resp.Body)
	fmt.Println(task.Name, resp.StatusCode, string(all))
	return nil
}

func expandArgs(origin string, parameters map[string]string) (string, []string) {
	// 缺少的参数
	var lessParameters []string
	return os.Expand(origin, func(p string) string {
		value := parameters[p]
		if value != "" {
			return value
		} else {
			lessParameters = append(lessParameters, p)
			return p
		}
	}), lessParameters
}

func getFormDataRequest(method string, url string, data string) (*http.Request, error) {
	var buffer bytes.Buffer

	writer := multipart.NewWriter(&buffer)
	for _, s := range strings.Split(data, ";") {
		pair := GetKeyValuePair(s, "=")
		if pair.Value[0] == '@' {
			filePath := strings.Replace(pair.Value, "@", "", 1)
			open, err := os.Open(filePath)
			if err != nil {
				return nil, errors.New("can not find file")
			}
			file, _ := writer.CreateFormFile(pair.Key, filepath.Base(filePath))
			io.Copy(file, open)
		} else {
			field, _ := writer.CreateFormField(pair.Key)
			field.Write([]byte(pair.Value))
		}
	}
	writer.Close()
	request, _ := http.NewRequest(method, url, &buffer)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request, nil
}

func getJsonRequest(method string, url string, data string) *http.Request {
	reader := strings.NewReader(data)
	request, _ := http.NewRequest(method, url, reader)
	return request
}

func getHeaders(headerKey string) map[string]string {
	headerMap := map[string]string{}
	headers := viper.GetStringSlice(headerKey)
	baseHeaders := LoadBaseConfig().Headers
	for _, s := range baseHeaders {
		pair := GetKeyValuePair(s, ":")
		headerMap[pair.Key] = pair.Value
	}

	for _, s := range headers {
		pair := GetKeyValuePair(s, ":")
		headerMap[pair.Key] = pair.Value
	}

	return headerMap
}

type KeyValue struct {
	Key   string
	Value string
}

func GetKeyValuePair(str string, s string) *KeyValue {
	splits := strings.Split(str, s)
	if len(splits) == 2 {
		return &KeyValue{splits[0], splits[1]}
	}
	return nil
}
