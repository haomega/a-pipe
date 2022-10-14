/*
run custom tasks as pipe (api request)

concept
- http request (json,file support)
- read and parse yml file (pipe-compose.yml)
- open and edit file?

feature
- command:
  - run pipe/task
  - list pipe/task

- pretty task running log
- env spec eg: ip,port

Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import "a-pipe/cmd"

func main() {
	cmd.Execute()
}
