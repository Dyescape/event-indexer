package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	home    = os.Getenv("HOME")
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "event-indexer",
		Short: "Event consumer to take Kafka events an index them into Elastic Search",
		Run: func(cmd *cobra.Command, _ []string) {
			if v, _ := cmd.Flags().GetBool("version"); v {
				cmd.Println(sprintVersion(cmd))
			} else {
				cmd.Help()
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("version", "v", false, "shows version information")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", `config file, by default searched in following paths:
- /etc/dyescape/event-indexer
- $HOME
- .`)
}

func initConfig() {
	viper.SetConfigName(".event-indexer")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AddConfigPath("/etc/dyescape/event-indexer")
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		} else {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}

// Execute runs the rule editor server
func Execute(v VersionInfo) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
