package task

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
)

type Pipe struct {
	Name  string
	Tasks []Task
}

func LoadPipe(pipeName string) (*Pipe, error) {
	pipeKey := "pipes." + pipeName

	var tasks []Task
	taskNames := viper.GetStringSlice(pipeKey)
	for _, taskName := range taskNames {
		task, err := LoadTask(taskName)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}
	return &Pipe{pipeName, tasks}, nil
}

func GetAppConfigPipes() []Pipe {
	var pipes []Pipe
	pipesMap := viper.GetStringMap("pipes")
	for pipeName := range pipesMap {
		pipe, err := LoadPipe(pipeName)
		if err != nil {
			panic(fmt.Errorf("load pipe error %w", err))
		}
		pipes = append(pipes, *pipe)
	}
	return pipes
}

func (pipe *Pipe) RunPipe() {
	bar := progressbar.Default(int64(len(pipe.Tasks)))
	for _, task := range pipe.Tasks {
		err := task.RequestApi()
		if err != nil {
			fmt.Println("Run task error "+task.Name, err)
		} else {
			bar.Add(1)
		}
	}
}
