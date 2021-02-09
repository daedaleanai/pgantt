package ptasks

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type Opts struct {
	Port           int      // Port to serve the on
	PhabricatorUri string   // Rendez-vous point with phabricator
	ApiKey         string   // Phabricator API key
	Projects       []string // List of projects to be handled
}

func NewOpts() (opts *Opts) {
	opts = new(Opts)
	opts.Port = 9999
	return
}

// Load the configuration data from a Yaml file
func (opts *Opts) LoadYaml(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("Unable to read the configuration file %s: %s", fileName, err)
	}

	err = yaml.Unmarshal(data, opts)
	if err != nil {
		return fmt.Errorf("Malformed config %s: %s", fileName, err)
	}

	return nil
}
