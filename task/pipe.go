package task

import (
	"fmt"
	"github.com/k0kubun/go-ansi"
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
	bar := progressbar.NewOptions(len(pipe.Tasks),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan]Run pipe's tasks...[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	for _, task := range pipe.Tasks {
		err := task.RequestApi()
		if err != nil {
			fmt.Println("Run task error "+task.Name, err)
		} else {
			bar.Add(1)
		}
	}
	bar.Finish()
}
