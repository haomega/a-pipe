/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"a-pipe/task"
	"fmt"
	"github.com/spf13/cobra"
)

var isRunTask bool
var taskArgs []string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run pipes/tasks",
	Long:  `run pipes, pipes are series tasks`,
	Run: func(cmd *cobra.Command, args []string) {
		if isRunTask {
			for _, taskName := range args {
				config, err := task.LoadTask(taskName)
				if err != nil {
					panic(fmt.Errorf("load task fail: %w", err))
				}
				fmt.Println("load task config", config)
				err = config.RequestApi(taskArgs)
				if err != nil {
					panic(fmt.Errorf("request for task fail: %w", err))
				}
			}
		} else {
			for _, pipeName := range args {
				pipe, err := task.LoadPipe(pipeName)
				if err != nil {
					panic(fmt.Errorf("load pipe error: %w", err))
				}
				pipe.RunPipe()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().BoolVarP(&isRunTask, "run task", "t", false, "run -t xxx")
	runCmd.PersistentFlags().StringArrayVarP(&taskArgs, "task args", "p", []string{}, "run -p 123 456")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
