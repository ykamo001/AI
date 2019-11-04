package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/ykamo001/backend/internal/paint"
	"github.com/ykamo001/backend/request"
	paintservice "github.com/ykamo001/backend/rpc/paint"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs the server",
	Long:  `Server for the AI service.`,
	Run:   runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	logger := setupLogger()
	router := mux.NewRouter().StrictSlash(false)

	paintProvider := paint.NewProvider(logger)
	paintServer := paintservice.NewPaintServer(paintProvider, nil)
	router.PathPrefix(paintServer.PathPrefix()).Handler(paintServer)

	err := http.ListenAndServe(":8080", request.WithRequestHeaders(router))
	logger.Error(err)
}
