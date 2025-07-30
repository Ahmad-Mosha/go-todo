package cmd

import (
	"fmt"
	"todo/models"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("completed", "c", false, "Show only completed tasks")
	listCmd.Flags().BoolP("pending", "p", false, "Show only pending tasks")
	listCmd.Flags().IntP("limit", "l", 0, "Limit number of tasks to show (0 = show all)")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all todo tasks",
	Run: func(cmd *cobra.Command, args []string) {

		showCompleted, _ := cmd.Flags().GetBool("completed")
		showPending, _ := cmd.Flags().GetBool("pending")
		limit, _ := cmd.Flags().GetInt("limit")

		tasks, err := csvStorage.ListTask()

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Printf("📝 No tasks found! Use 'todo add \"task\"' to create one.\n")
			return
		}
		filteredTasks := filterTasks(tasks, showCompleted, showPending)
		if limit > 0 && len(filteredTasks) > limit {
			filteredTasks = filteredTasks[:limit]
		}

		if len(filteredTasks) == 0 {
			fmt.Printf("📝 No tasks match your filter criteria.\n")
			return
		}

		filterInfo := ""
		if showCompleted {
			filterInfo = " (completed only)"
		} else if showPending {
			filterInfo = " (pending only)"
		}

		fmt.Printf("\n📋 Your Tasks (%d total%s):\n\n", len(filteredTasks), filterInfo)

		for _, task := range filteredTasks {
			status := "⏳ Pending"
			if task.Completed {
				status = "✅ Done"
			}
			fmt.Printf("[%d] %s - %s\n", task.ID, task.Description, status)
		}
		fmt.Printf("\n")
	},
}

func filterTasks(tasks []*models.Task, showCompleted, showPending bool) []*models.Task {
	// If no filter flags are set, return all tasks
	if !showCompleted && !showPending {
		return tasks
	}

	var filtered []*models.Task
	for _, task := range tasks {
		if showCompleted && task.Completed {
			filtered = append(filtered, task)
		} else if showPending && !task.Completed {
			filtered = append(filtered, task)
		}
	}
	return filtered
}
