package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	str "github.com/luigi-riefolo/eGO/l10n/en_EN"
)

// TODO:
// make the list of services dynamic
// support slices

//ALFA_SERVER_PORT=4444 go run main.go service -opts CONFIG=../../../conf/global_conf.toml,ALFA_SERVER_PORT=90932

const confKey = "CONFIG_FILE"

// list of overwritable config values
// NOTE: the order is config, environment
// variable and command-line argument
var (
	conf = Config{}

	cmdArgMap = make(map[string]string)
	confMap   = make(map[string]*reflect.Value)
)

func init() {
	parseArgs()
}

// PrintJSONConfig prints the services configuration in JSON format.
func PrintJSONConfig() error {
	conf, err := GetConfig()
	if err != nil {
		return fmt.Errorf("Config file could not be parsed: %v", err)
	}

	buf, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("Configuration could not be marshalled: %v", err)
	}
	fmt.Printf("%s\n", string(buf))

	return nil
}

// GetConfig returns a Config struct representing a TOML file.
func GetConfig() (Config, error) {

	// set the config file path
	conf.ConfigFile = confFile
	if val, ok := os.LookupEnv(confKey); ok {
		conf.ConfigFile = val
	}

	// collect the command-line arguments
	conf.processArgs()
	if val, ok := cmdArgMap[confKey]; ok {
		conf.ConfigFile = val
	}

	if conf.ConfigFile == "" {
		return conf, fmt.Errorf(str.NoConfigFile)
	}

	var err error
	conf.ConfigFile, err = filepath.Abs(conf.ConfigFile)
	if err != nil {
		return conf, err
	}

	if _, err = toml.DecodeFile(conf.ConfigFile, &conf); err != nil {
		return conf, err
	}

	if err := conf.processConf(); err != nil {
		return conf, err
	}

	return conf, nil
}

// processConf populates the configuration with each respective value.
func (conf *Config) processConf() error {

	root := reflect.ValueOf(conf).Elem()

	// map each configuration field
	traverse(&root, "")

	for k, v := range confMap {

		key := strings.ToUpper(k)

		// overwrite config values with the
		// respective environment variable
		val, _ := os.LookupEnv(key)

		if argVal, ok := cmdArgMap[key]; ok {
			val = argVal
		}

		if err := overwrite(v, val); err != nil {
			return err
		}
	}

	return nil
}

// overwrite sets a configuration value.
func overwrite(arg *reflect.Value, argVal string) error {
	if argVal == "" {
		return nil
	}

	switch arg.Kind() {

	//case reflect.Slice:

	case reflect.Int:

		n, err := strconv.Atoi(argVal)
		if err != nil {
			return err
		}

		if arg.IsValid() && arg.CanSet() {
			arg.Set(reflect.ValueOf(n))
		}

	case reflect.Bool:
		fallthrough

	case reflect.String:

		if arg.IsValid() && arg.CanSet() {
			arg.Set(reflect.ValueOf(argVal))
		}
	}

	return nil
}

// traverse recursively traverses the configuration and maps each field.
func traverse(root *reflect.Value, path string) {

	if root.Kind() != reflect.Struct {
		path = strings.ToLower(path)
		confMap[path] = root

		return
	}

	sep := "_"
	for i := 0; i < reflect.Indirect(*root).NumField(); i++ {

		field := reflect.Indirect(*root).Field(i)

		fieldName := reflect.Indirect(*root).Type().Field(i).Name

		if path == "" {
			sep = ""
		}

		newPath := fmt.Sprintf("%s%s%s", path, sep, fieldName)

		traverse(&field, newPath)
	}
}

// processArgs processes the list of comma-separated command-line arguments.
func (conf *Config) processArgs() {

	if opts == "" {
		return
	}

	optList := strings.Split(opts, ",")

	for _, op := range optList {
		arg := strings.Split(op, "=")
		cmdArgMap[arg[0]] = arg[1]
	}
}
