package cmd

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"adiachenko/go-scaffold/routes"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start http server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Server is ready to handle requests at port 8000")
		log.Fatal(http.ListenAndServe(":8000", routes.Register()))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
