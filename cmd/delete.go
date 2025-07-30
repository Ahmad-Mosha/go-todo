package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Add flags
	deleteCmd.Flags().BoolP("force", "f", false, "Force delete without confirmation")
}

var deleteCmd = &cobra.Command{
	Use:   "delete [task ID]",
	Short: "Delete a todo task",
	Long:  "Delete a task by its ID number",
	Args:  cobra.ExactArgs(1), // Requires exactly 1 argument (the ID)
	Run: func(cmd *cobra.Command, args []string) {
		// Parse the task ID
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("❌ Error: Invalid task ID '%s'. Please provide a valid number.\n", args[0])
			return
		}

		// Get force flag
		force, _ := cmd.Flags().GetBool("force")

		// If not forced, ask for confirmation
		if !force {
			fmt.Printf("⚠️  Are you sure you want to delete task [%d]? (y/N): ", id)
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
				fmt.Printf("❌ Delete cancelled.\n")
				return
			}
		}

		// Delete the task
		err = csvStorage.DeleteTask(id)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		// Success message
		fmt.Printf("✅ Task [%d] deleted successfully!\n", id)
		fmt.Printf("💡 Use 'todo list' to see your remaining tasks\n")
	},
}
