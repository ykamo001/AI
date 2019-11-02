package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	//Version is injected by go (should be a tag name)
	Version = "None"
	// Buildstamp is a timestamp (injected by go) of the build time
	Buildstamp = "None"
	// Githash is the tag for current hash the build represents
	Githash = "None"
	host    = "None"
)
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ai",
	Short: "Service that runs all internal AI projects",
	Long:  `Service holds an AI project that finds best features for classification.`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.IsSet("SUBCOMMAND") {
			switch viper.GetString("SUBCOMMAND") {
			case "server":
				runServer(cmd, args)
			}
		}
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	viper.AutomaticEnv()

	defaults := map[string]interface{}{
		"PORT": 8080,
	}

	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	var err error
	host, err = os.Hostname()
	if err != nil {
		logrus.Panicln("unable to get Hostname", err)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.WithFields(logrus.Fields{
		"Version":   Version,
		"BuildTime": Buildstamp,
		"Bithash":   Githash,
		"Host":      host,
	}).Info("Service Statup")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ai.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ai" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ai")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func setupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)

	return logger
}
