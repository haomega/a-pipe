/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var aPipeCompose string
var rootCmd = &cobra.Command{
	Use:   "a-pipe",
	Short: "run custom tasks as pipe",
	Long:  `run custom tasks as pipe`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.a-pipe.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// a-pipe-compose.yml 文件
	home, err := homedir.Dir()
	defaultConfigFile := home + "/.a-pipe/a-pipe-compose.yml"
	//if err != nil {
	//	fmt.Errorf("can't dedct home dir")
	//}
	rootCmd.PersistentFlags().StringVar(&aPipeCompose, "pipe config compose file", defaultConfigFile, "pipe config compose file")
	viper.SetConfigName("a-pipe-compose.yml") // name of config file (without extension)
	viper.SetConfigType("yaml")               // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/.a-pipe")      // call multiple times to add many search paths
	viper.AddConfigPath(".")                  // optionally look for config in the working directory
	err = viper.ReadInConfig()                // Find and read the config file
	if err != nil {                           // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = viper.WriteConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("write error config file: %w", err))
	}
	used := viper.ConfigFileUsed()
	fmt.Println(used)
	getString := viper.GetString("a-pipe.config.domain")
	fmt.Println(getString)
}
