package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all todo tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := csvStorage.ListTask()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Printf("📝 No tasks found! Use 'todo add \"task\"' to create one.\n")
			return
		}

		fmt.Printf("\n📋 Your Tasks (%d total):\n\n", len(tasks))
		for _, task := range tasks {
			status := "⏳ Pending"
			if task.Completed {
				status = "✅ Done"
			}
			fmt.Printf("[%d] %s - %s\n", task.ID, task.Description, status)
		}
		fmt.Printf("\n")
	},
}
