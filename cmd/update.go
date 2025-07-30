package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)

	// Add flags
	updateCmd.Flags().StringP("description", "d", "", "Update task description")
	updateCmd.Flags().BoolP("complete", "c", false, "Mark task as completed")
	updateCmd.Flags().BoolP("incomplete", "i", false, "Mark task as incomplete")
}

var updateCmd = &cobra.Command{
	Use:   "update [task ID]",
	Short: "Update a todo task",
	Long:  "Update a task's description or completion status",
	Args:  cobra.ExactArgs(1), // Requires exactly 1 argument (the ID)
	Run: func(cmd *cobra.Command, args []string) {
		// Parse the task ID
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("❌ Error: Invalid task ID '%s'. Please provide a valid number.\n", args[0])
			return
		}

		// Get flag values
		newDescription, _ := cmd.Flags().GetString("description")
		markComplete, _ := cmd.Flags().GetBool("complete")
		markIncomplete, _ := cmd.Flags().GetBool("incomplete")

		// Validate flags
		if markComplete && markIncomplete {
			fmt.Printf("❌ Error: Cannot use both --complete and --incomplete flags together.\n")
			return
		}

		if newDescription == "" && !markComplete && !markIncomplete {
			fmt.Printf("❌ Error: Please specify what to update using flags:\n")
			fmt.Printf("   --description \"new description\" or -d \"new description\"\n")
			fmt.Printf("   --complete or -c (mark as completed)\n")
			fmt.Printf("   --incomplete or -i (mark as incomplete)\n")
			return
		}

		// Prepare update parameters
		var descPtr *string
		var completedPtr *bool

		if newDescription != "" {
			descPtr = &newDescription
		}

		if markComplete {
			completed := true
			completedPtr = &completed
		} else if markIncomplete {
			completed := false
			completedPtr = &completed
		}

		// Update the task
		updatedTask, err := csvStorage.UpdateTask(id, descPtr, completedPtr)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		// Success message with updated task details
		fmt.Printf("\n✅ Task updated successfully!\n")
		fmt.Printf("┌─────────────────────────────────────────┐\n")
		fmt.Printf("│ ID: %-3d                               │\n", updatedTask.ID)
		fmt.Printf("│ Task: %-30s │\n", truncateString(updatedTask.Description, 30))

		status := "Pending"
		if updatedTask.Completed {
			status = "Completed"
		}
		fmt.Printf("│ Status: %-27s │\n", status)
		fmt.Printf("│ Updated: %-26s │\n", updatedTask.CreatedAt.Format("2006-01-02 15:04"))
		fmt.Printf("└─────────────────────────────────────────┘\n")
		fmt.Printf("\n💡 Use 'todo list' to see all your tasks\n")
	},
}

// Helper function to truncate long descriptions
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
