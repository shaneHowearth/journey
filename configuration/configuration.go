// Package configuration -
package configuration

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/kabukky/journey/filenames"
)

// Configuration - settings that are neccesary for server configuration
type Configuration struct {
	HTTPHostAndPort  string
	HTTPSHostAndPort string
	HTTPSUsage       string
	URL              string
	HTTPSURL         string
	UseLetsEncrypt   bool
}

// NewConfiguration -
func NewConfiguration() *Configuration {
	var config Configuration
	err := config.load()
	if err != nil {
		log.Println("Warning: couldn't load " + filenames.ConfigFilename + ", creating new config file.")
		err = config.create()
		if err != nil {
			log.Fatal("Fatal error: Couldn't create configuration.")
			return nil
		}
		err = config.load()
		if err != nil {
			log.Fatal("Fatal error: Couldn't load configuration.")
			return nil
		}
	}
	return &config
}

// Config - Global config - thread safe and accessible from all packages
var Config = NewConfiguration()

func (c *Configuration) save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filenames.ConfigFilename, data, 0600)
}

func (c *Configuration) load() error {
	configWasChanged := false
	data, err := ioutil.ReadFile(filenames.ConfigFilename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}
	// Make sure the url is in the right format
	// Make sure there is no trailing slash at the end of the url
	if strings.HasSuffix(c.URL, "/") {
		c.URL = c.URL[0 : len(c.URL)-1]
		configWasChanged = true
	}
	if !strings.HasPrefix(c.URL, "http://") && !strings.HasPrefix(c.URL, "https://") {
		c.URL = "http://" + c.URL
		configWasChanged = true
	}
	// Make sure the https url is in the right format
	// Make sure there is no trailing slash at the end of the https url
	if strings.HasSuffix(c.HTTPSURL, "/") {
		c.HTTPSURL = c.HTTPSURL[0 : len(c.HTTPSURL)-1]
		configWasChanged = true
	}
	if strings.HasPrefix(c.HTTPSURL, "http://") {
		c.HTTPSURL = strings.Replace(c.HTTPSURL, "http://", "https://", 1)
		configWasChanged = true
	} else if !strings.HasPrefix(c.HTTPSURL, "https://") {
		c.HTTPSURL = "https://" + c.HTTPSURL
		configWasChanged = true
	}
	// Make sure there is no trailing slash at the end of the url
	if strings.HasSuffix(c.HTTPSURL, "/") {
		c.HTTPSURL = c.HTTPSURL[0 : len(c.HTTPSURL)-1]
		configWasChanged = true
	}
	// Check if all fields are filled out
	cReflected := reflect.ValueOf(*c)
	for i := 0; i < cReflected.NumField(); i++ {
		if cReflected.Field(i).Interface() == "" {
			log.Println("Error: " + filenames.ConfigFilename + " is corrupted. Did you fill out all of the fields?")
			return errors.New("configuration corrupted")
		}
	}
	// Save the changed config
	if configWasChanged {
		err = c.save()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Configuration) create() error {
	// TODO: Change default port
	c = &Configuration{HTTPHostAndPort: ":8084", HTTPSHostAndPort: ":8085", HTTPSUsage: "None", URL: "127.0.0.1:8084", HTTPSURL: "127.0.0.1:8085"}
	err := c.save()
	if err != nil {
		log.Println("Error: couldn't create " + filenames.ConfigFilename)
		return err
	}

	return nil
}
