package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Task struct {
	ID          int
	Name        string
	Description string
	Completed   bool
}

// Func get next ID
func getNextID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID
}

// Function to read the json file
func readJSONFile(filePath string) ([]Task, error) {
	var tasks []Task
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &tasks); err != nil {
			return nil, err
		}
	}
	return tasks, nil
}

func appendToJSONFile(filePath string) error {
	tasks, err := readJSONFile(filePath)
	if err != nil {
		fmt.Printf("%s", err)
	}
	//Get the nextID
	nextID := getNextID(tasks) + 1

	fmt.Println("Add Task")
	var name string
	var description string
	fmt.Scanln(&name)
	fmt.Scanln(&description)
	newTask := Task{
		ID:          nextID,
		Name:        name,
		Description: description,
		Completed:   false,
	}
	tasks = append(tasks, newTask)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(tasks)
}

func deleteTheTask(filePath string) error {
	tasks, err := readJSONFile(filePath)
	if err != nil {
		return err
	}
	fmt.Println("Enter ID to delete Task")
	var toDeleteID int
	fmt.Scanln(&toDeleteID)
	if toDeleteID >= 0 && toDeleteID < len(tasks) {
		tasks = append(tasks[:toDeleteID], tasks[toDeleteID+1:]...)
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(tasks)

}

func makeItComplete(filePath string) error {
	fmt.Println("Completed What?")
	var toChange int
	fmt.Scanln(&toChange)
	tasks, err := readJSONFile(filePath)
	if err != nil {
		fmt.Printf("%s", err)
	}
	for index := range tasks {
		if index == toChange {
			tasks[index].Completed = !tasks[index].Completed
			break
		}
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(tasks)
}

func main() {

	//Filepath
	filePath := "tasks.json"
	_, err := os.Stat(filePath)

	//If file exists then fine or create the file
	if os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file")
			return
		}
		defer file.Close()
		fmt.Println("File 'tasks.json' has been created.")
	} else if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Printf("")

	//Creating a infinite loop for user input
	for {
		fmt.Printf("1. Add Task\n2.Delete Task\n3.Complete the task Yo!\n4.List All the tasks\n5.Exit\n")
		var choice int
		toExit := false
		fmt.Println("Enter your choice:")
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			//Add the task
			err := appendToJSONFile(filePath)
			if err != nil {
				fmt.Printf("%s", err)
			}

		case 2:
			//Delete the task
			deleteTheTask(filePath)

		case 3:
			//Complete the task
			makeItComplete(filePath)

		case 4:
			//List the tasks
			jsonData, err := readJSONFile(filePath)
			if err != nil {
				fmt.Printf("%s", err)
			}
			for index, task := range jsonData {
				fmt.Printf("%d\t%-15s\t%-25s\t%t\n", index, task.Name, task.Description, task.Completed)
			}

		case 5:
			fmt.Println("Breaking out of the loop")
			toExit = true
		}
		if toExit {
			break
		}
	}
	fmt.Printf("Done!\n")
}
