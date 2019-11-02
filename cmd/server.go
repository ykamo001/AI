package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/ykamo001/ai/internal/featureselection"
	featureselectionservice "github.com/ykamo001/ai/rpc/featureselection"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs the server",
	Long:  `Server for the ai service.`,
	Run:   runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	logger := setupLogger()

	router := mux.NewRouter().StrictSlash(true)

	featureSelectionProvider := featureselection.NewProvider(logger)
	featureSelectionServer := featureselectionservice.NewFeatureSelectionServer(featureSelectionProvider, nil)
	router.PathPrefix(featureSelectionServer.PathPrefix()).Handler(featureSelectionServer)

	err := http.ListenAndServe(":8080", router)
	logger.Error(err)
}
