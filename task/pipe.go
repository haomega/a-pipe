package task

import (
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
	"strings"
)

type Pipe struct {
	Name  string
	Tasks []TaskWithArg
}

type TaskWithArg struct {
	Task       Task
	Parameters map[string]string
}

func LoadPipe(pipeName string) (*Pipe, error) {
	pipeKey := "pipes." + pipeName

	var tasks []TaskWithArg
	taskDefList := viper.GetStringSlice(pipeKey)
	for _, taskDef := range taskDefList {
		var taskName string
		parameters := map[string]string{}
		splits := strings.Split(taskDef, " ")
		if len(splits) > 1 {
			taskName = splits[0]
			for i := 1; i < len(splits); i++ {
				pair := GetKeyValuePair(splits[i], "=")
				parameters[pair.Key] = pair.Value
			}
		} else {
			taskName = taskDef
		}
		task, err := LoadTask(taskName)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, TaskWithArg{Task: *task, Parameters: parameters})
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

func (pipe *Pipe) RunPipe(parameters map[string]string) {
	bar := progressbar.NewOptions(len(pipe.Tasks),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan]Run...[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	for _, taskWithArgs := range pipe.Tasks {
		task := taskWithArgs.Task
		ps := taskWithArgs.Parameters
		// merge parameters
		for p := range parameters {
			ps[p] = parameters[p]
		}
		bar.Add(1)
		err := task.RequestApi(ps)
		if err != nil {
			fmt.Println("Run task error "+task.Name, err)
		}
	}
	bar.Finish()
}
