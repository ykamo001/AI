package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/ykamo001/ai/internal/eightpuzzle"
	"github.com/ykamo001/ai/request"
	eightpuzzleservice "github.com/ykamo001/ai/rpc/eightpuzzle"
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

	eightPuzzleProvider := eightpuzzle.NewProvider(logger)
	eightPuzzleServer := eightpuzzleservice.NewEightPuzzleServer(eightPuzzleProvider, nil)
	router.PathPrefix(eightPuzzleServer.PathPrefix()).Handler(eightPuzzleServer)

	err := http.ListenAndServe(":8080", request.WithRequestHeaders(router))
	logger.Error(err)
}
