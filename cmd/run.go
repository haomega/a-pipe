/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"a-pipe/task"
	"fmt"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run pipes",
	Long:  `run pipes, pipes are series tasks`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
		fmt.Println(args)
	},
}

var runTask = &cobra.Command{
	Use:   "task",
	Short: "t",
	Long:  "run tasks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run tasks called")
		for _, taskName := range args {
			taskConfig := task.LoadTaskConfig(taskName)
			fmt.Println(taskConfig)
			task.HttpRequest(&taskConfig)
		}
	},
}

func init() {
	runCmd.AddCommand(runTask)
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
