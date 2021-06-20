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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "./conf.d/env.yaml", "config file")
	rootCmd.PersistentFlags().StringVarP(&gitCommitNum, "version", "v", "unknown", "git commit hash")
	rootCmd.PersistentFlags().StringVarP(&buildTime, "buildTime", "b", time.Now().String(), "binary build time")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	logger.Init("debug")

	viper.SetConfigName("env")                             // name of config file (without extension)
	viper.SetConfigType("yaml")                            // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./conf.d")                        // path to look for the config file in
	viper.AutomaticEnv()                                   //
	viper.SetEnvPrefix("DongTzu")                          //
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) //

	// 有指定 config file 時直接讀取, 不透過已設定的 config path 尋找
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("ReadInConfig file failed: %v", err)
		setDefaultConfig()
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
		return err
	}

	config.SetConfig(c)

	return nil
}

func setDefaultConfig() {
	logger.Debugf("Set default config.")

	// ArangoDB
	viper.SetDefault("arangoDB.addr", "http://127.0.0.1:8529/")
	viper.SetDefault("arangoDB.dbName", "_system")
	viper.SetDefault("arangoDB.username", "root")
	viper.SetDefault("arangoDB.password", "")
	viper.SetDefault("arangoDB.connLimit", 40)

	// Fiber
	viper.SetDefault("fiber.port", ":4000")

	// Zoom
	viper.SetDefault("zoom.meetingExtendedTime", 10)

	// Github
	viper.SetDefault("github.apiToken", "")
	viper.SetDefault("github.repoURL", "")
}
