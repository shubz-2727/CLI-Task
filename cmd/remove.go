package cmd

import (
	"fmt"
	"strconv"
	"task/db"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a task from task list.",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		if len(args) == 0 {
			fmt.Println("No argument passed..")
		}
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to passed the argument ", arg)
			} else {
				ids = append(ids, id)
			}
		}

		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}

		for _, id := range ids {

			if id <= 0 || id > len(tasks) {
				fmt.Println("Invalid task number ", id)
				continue
			}
			task := tasks[id-1]
			err := db.DeleteTask(task.Key)
			if err != nil {
				fmt.Println("failed to remove task ", id, err)
			} else {
				fmt.Printf("Remove \"%s\" from your task list.\n", task.Value)

			}

		}
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)
}
