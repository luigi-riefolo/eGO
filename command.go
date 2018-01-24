package main

import (
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	service string
	version string
	build   string
	rootCmd = &cobra.Command{}
)

/*
{
    "host": {
        "address": "localhost",
        "port": 5799
    },
    "datastore": {
        "metric": {
            "host": "127.0.0.1",
            "port": 3099
        },
        "warehouse": {
            "host": "198.0.0.1",
            "port": 2112
        }
    }
}

type config struct {
	Port int
	Name string
	PathMap string `mapstructure:"path_map"`
}

var C config

err := Unmarshal(&C)
if err != nil {
	t.Fatalf("unable to decode into struct, %v", err)
}
*/
type conf struct{}

func init() {
	/*
		viper.SetDefault("ContentDir", "content")
		viper.SetConfigName("config")         // name of config file (without extension)
		viper.AddConfigPath("./conf/")        // path to look for the config file in
		viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
		viper.AddConfigPath(".")              // optionally look for config in the working directory
		err := viper.ReadInConfig()           // Find and read the config file
		if err != nil {                       // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
		// Watching and re-reading config files
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("config file changed:", e.Name)
		})
		// alternatively, you can create a new viper instance.
		var runtimeViper = viper.New()

		err = viper.AddSecureRemoteProvider(
			"etcd",
			"http://127.0.0.1:4001",
			"/config/hugo.json",
			"/etc/secrets/mykeyring.gpg")

		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
		runtimeViper.SetConfigType("yaml") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop"

		// read from remote config the first time.
		err = runtimeViper.ReadRemoteConfig()

		// unmarshal config
		runtimeConf := &conf{}
		err = runtimeViper.Unmarshal(runtimeConf)
		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %s", err))
		}

		// open a goroutine to watch remote changes forever
		go func() {
			tt := time.NewTicker(time.Minute)
			stop := make(chan struct{})
			for {
				select {
				case <-tt.C:

					// currently, only tested with etcd support
					if err := runtimeViper.WatchRemoteConfig(); err != nil {
						log.Printf("unable to read remote config: %v", err)
						continue
					}

					// unmarshal new config into our runtime config struct.
					// you can also use channel to implement a signal to notify
					// the system of the changes
					if err := runtimeViper.Unmarshal(&runtimeConf); err != nil {
						// Handle errors reading the config file
						panic(fmt.Errorf("fatal error config file: %s", err))
					}
				case <-stop:
					tt.Stop()
				}
			}
		}()
	*/
	_ = viper.GetString("datastore.metric.host")

	///////////////////////////////////////////////////////////////
	// cobra
	rootCmd = &cobra.Command{
		Use:   service,
		Short: "short desc",
		Long:  "long desc",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("OK")
		},
	}

	var startCmd = &cobra.Command{
		Use:     "start",
		Short:   fmt.Sprintf("Start the %s service", service),
		Long:    fmt.Sprintf("Start the %s service", service),
		Example: fmt.Sprintf("%s start", service),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting the %s service", service)
		},
	}
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print the version number of %s service", service),
		Long:  fmt.Sprintf("Print the version number of %s service", service),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version number of %s service is %s\n", service, version)
		},
	}
	/*
		cobra.OnInitialize(initConfig)
		rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")

		rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
		rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
		rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
		rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
		err := viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
		err = viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
		err = viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
		viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
		viper.SetDefault("license", "apache")
	*/
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)

	/*
		// alternatively, you can create a new viper instance.
		var runtimeViper = viper.New()
		err := runtimeViper.AddRemoteProvider("etcd", "http://127.0.0.1:4001", "/config/hugo.yml")
		if err != nil {
			log.Fatal(err)
		}*/
}

var (
	cfgFile, projectBase string
	_                    = projectBase
)

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
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

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
func main() {
	log.Printf("Starting '%s' service, version '%s' (%s)\n", service, version, build)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	//log.Printf("Server listening on %s\n", port)
	//log.Fatalln(http.ListenAndServe(port, nil))
}
