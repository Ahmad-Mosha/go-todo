package cmd

import (
	"fmt"
	"strings"
	"todo/internal/storage"

	"github.com/spf13/cobra"
)

var csvStorage = storage.NewCSVStorage("todos.csv")

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "add task",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		description := strings.Join(args, " ")

		task, err := csvStorage.AddTask(description)
		if err != nil {
			fmt.Printf("error adding task: %v\n", err)
			return
		}
		fmt.Printf("Task added: [%d] %s\n", task.ID, task.Description)
	},
}
