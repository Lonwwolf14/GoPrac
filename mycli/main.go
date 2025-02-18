package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type Task struct {
	Title string
	Done  bool
}

var taskFile = "tasks.json"

func loadTasks() {
	data, err := ioutil.ReadFile(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			TaskList = []Task{}
			return
		}
		fmt.Printf("Error reading the file %s", err)
		return
	}
	err = json.Unmarshal(data, &TaskList)
	if err != nil {
		fmt.Println("Error unmarshalling tasks:", err)
	}
}

func saveTasks() {
	data, err := json.MarshalIndent(TaskList, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling tasks:", err)
		return
	}

	err = ioutil.WriteFile(taskFile, data, 0644)
	if err != nil {
		fmt.Println("Error writing task file:", err)
	}
}

var TaskList []Task

var rootCmd = &cobra.Command{
	Use:   "taskcli",
	Short: "A simple task management CLI",
	Long:  `A CLI tool for managing tasks. You can add, list, and remove tasks from your task list.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		loadTasks()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		saveTasks()
	},
}

var addCmd = &cobra.Command{
	Use:   "add [task]",
	Short: "Add a new task",
	Long:  "Append your new task to the once already you haven't completed Bruhh!!",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		task := Task{Title: args[0], Done: false}
		TaskList = append(TaskList, task)
		fmt.Printf("Task has been added : %s\n", task.Title)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  "List all tasks in the task list",
	Run: func(cmd *cobra.Command, args []string) {
		if len(TaskList) == 0 {
			fmt.Println("No tasks in the list.")
			return
		}
		for i, task := range TaskList {
			status := "Not Done"
			if task.Done {
				status = "Done"
			}
			fmt.Printf("%d. %s - %s\n", i+1, task.Title, status)
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove [index]",
	Short: "Remove a task by index",
	Long:  "Remove a task from the task list by specifying its index",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := strconv.Atoi(args[0])
		if err != nil || index < 1 || index > len(TaskList) {
			fmt.Println("Invalid index provided")
			return
		}
		TaskList = append(TaskList[:index-1], TaskList[index:]...)
		fmt.Printf("Removed task at index %d\n", index)
	},
}

// doneCmd represents the mark done subcommand
var doneCmd = &cobra.Command{
	Use:   "done [index]",
	Short: "Mark a task as done by index",
	Long:  "Mark a task as done from the task list by specifying its index",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := strconv.Atoi(args[0])
		if err != nil || index < 1 || index > len(TaskList) {
			fmt.Println("Invalid index provided")
			return
		}
		TaskList[index-1].Done = true
		fmt.Printf("Marked task at index %d as done\n", index)
	},
}

func init() {
	rootCmd.AddCommand(addCmd, listCmd, removeCmd, doneCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
