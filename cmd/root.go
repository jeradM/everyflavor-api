package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var rootCmd = &cobra.Command{
	Use: "run",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file location")

	migrateCmd.PersistentFlags().IntP("version", "v", -1, "Migration version target")
	viper.BindPFlag("migrationVersion", migrateCmd.PersistentFlags().Lookup("version"))

	testDataCmd.PersistentFlags().IntP("num-recipes", "n", 500, "Number of recipes to generate")
	viper.BindPFlag("numRecipes", testDataCmd.PersistentFlags().Lookup("num-recipes"))

	rootCmd.AddCommand(appCmd, migrateCmd, testDataCmd)
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in current and home directories with name ".everyflavor.yaml".
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("everyflavor")
	}

	viper.AutomaticEnv()
	viper.SetDefault("serverAddr", "0.0.0.0:8099")
	viper.SetDefault("redisUrl", "localhost:6379")
	viper.SetDefault("redisDb", 0)
	viper.SetDefault("redisAuthKey", []string{"a67ee4161400c95ed5e880aca174beb240ed666a0d9be9b942813317f144891f"})
	viper.SetDefault("showSql", false)
	viper.SetDefault("corsAllowedOrigins", nil)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
