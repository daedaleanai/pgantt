package ptasks

import (
	"fmt"
	"reflect"

	"github.com/thought-machine/gonduit"
	"github.com/thought-machine/gonduit/core"
	"github.com/thought-machine/gonduit/requests"

	log "github.com/sirupsen/logrus"
)

type Phabricator struct {
	c *gonduit.Conn
}

type Task struct {
}

func (p *Phabricator) getProjectPhid(name string) (string, error) {
	req := requests.ProjectQueryRequest{Names: []string{"Platforms"}}
	res, err := p.c.ProjectQuery(req)
	if err != nil {
		return "", err
	}

	keys := reflect.ValueOf(res.Data).MapKeys()

	if len(keys) == 0 {
		return "", fmt.Errorf("Project not found: %s", name)
	}

	return keys[0].Interface().(string), err
}

func NewPhabricator(endpoint, key string) (*Phabricator, error) {
	conn, err := gonduit.Dial(endpoint, &core.ClientOptions{APIToken: key})
	if err != nil {
		return nil, err
	}

	return &Phabricator{conn}, nil
}

func (p *Phabricator) GetTasksForProject(name string) ([]Task, error) {
	phid, err := p.getProjectPhid(name)
	if err != nil {
		return nil, err
	}
	log.Errorf("%v", phid)
	return nil, nil
}
