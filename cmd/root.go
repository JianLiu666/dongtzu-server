package cmd

import (
	"dongtzu/config"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.geax.io/demeter/gologger/logger"
)

var cfgFile string
var gitCommitNum string
var buildTime string

var rootCmd = &cobra.Command{
	Use:               "root",
	Short:             "choose which one command to run: server",
	Long:              ``,
	PersistentPreRunE: PersistentPreRunBeforeCommandStartUp,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "./conf.d/env.yaml", "config file")
	rootCmd.PersistentFlags().StringVarP(&gitCommitNum, "version", "v", "unknown", "git commit hash")
	rootCmd.PersistentFlags().StringVarP(&buildTime, "buildTime", "b", time.Now().String(), "binary build time")

	logger.Init("debug")
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("DongTzu")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("ReadInConfig file failed: %v", err)
	} else {
		logger.Debugf("Using config file: %v", viper.ConfigFileUsed())
	}
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func PersistentPreRunBeforeCommandStartUp(cmd *cobra.Command, args []string) error {
	goVersion := runtime.Version()
	osName := runtime.GOOS
	architecture := runtime.GOARCH
	logger.Debugf("============================================================")
	logger.Debugf("Build on %s", buildTime)
	logger.Debugf("GoVersion: %s", goVersion)
	logger.Debugf("GitCommitNum: %s", gitCommitNum)
	logger.Debugf("OS: %s", osName)
	logger.Debugf("Architecture: %s", architecture)
	logger.Debugf("============================================================")

	c, err := config.NewFromViper()
	if err != nil {
		logger.Errorf("Init config failed: %v", err)
	} else {
		config.SetConfig(c)
	}

	return nil
}
