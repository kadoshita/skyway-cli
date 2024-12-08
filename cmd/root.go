package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "skyway-cli",
	Short: "A CLI tool for SkyWay developers",
	Long: `skyway-cli is a command-line interface tool designed to support the development of SkyWay.
By using this tool, you can automate and streamline the development process involving SkyWay's API.
Its main features include easy API calls, and token generate.
For engineers developing applications with SkyWay, this tool contributes to project efficiency and quality improvement.`,
	Version: "0.0.1",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func GenDocs() error {
	if _, err := os.Stat("./docs"); os.IsNotExist(err) {
		if err := os.Mkdir("./docs", 0755); err != nil {
			return err
		}
	}
	return doc.GenMarkdownTree(rootCmd, "./docs")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.skyway-cli.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".skyway" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".skyway-cli")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Debug("Using config", "file", viper.ConfigFileUsed())
	}
}
