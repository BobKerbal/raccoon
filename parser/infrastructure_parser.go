package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
)

//InfrastructureFileParser is the type that is in charge of parsing the files
//with the infrastructure information about a Raccoon job.
type InfrastructureFileParser struct {
	FileParser
}

//Build returns an Infrastructure from a io.Reader
func (t *InfrastructureFileParser) Build(r io.Reader) (*raccoon.Infrastructure, error) {
	var infrastructure raccoon.Infrastructure
	err := json.NewDecoder(r).Decode(&infrastructure)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON: " + err.Error())
	}

	t.takeAuthAtClusterLevel(&infrastructure)

	return t.checkErrors(&infrastructure)
}

//takeAuthAtClusterLevel will check if the user has written the authentication of a cluster in the
//cluster definition instead of on each individual host. In such case it will take the information
//and inject it on each host
func (t *InfrastructureFileParser) takeAuthAtClusterLevel(i *raccoon.Infrastructure) error {
	for k, cluster := range i.Infrastructure {
		var username, authFile, password string
		var interactive bool

		if i.Infrastructure[k].Username != "" {
			username = i.Infrastructure[k].Username
		}

		if i.Infrastructure[k].Password != "" {
			password = i.Infrastructure[k].Password
		}

		if i.Infrastructure[k].IdentityFile != "" {
			authFile = i.Infrastructure[k].IdentityFile
		}

		if i.Infrastructure[k].InteractiveAuth {
			interactive = i.Infrastructure[k].InteractiveAuth
		}

		for j := range cluster.Hosts {
			if i.Infrastructure[k].Hosts[j].Username == "" {
				i.Infrastructure[k].Hosts[j].Username = username
			}

			if i.Infrastructure[k].Hosts[j].Password == "" {
				i.Infrastructure[k].Hosts[j].Password = password
			}

			if i.Infrastructure[k].Hosts[j].IdentityFile == "" {
				i.Infrastructure[k].Hosts[j].IdentityFile = authFile
			}

			if !i.Infrastructure[k].Hosts[j].InteractiveAuth {
				i.Infrastructure[k].Hosts[j].InteractiveAuth = interactive
			}
		}
	}

	return nil
}

//checkErrors is used to perform error checking on Infrastructure json file
func (t *InfrastructureFileParser) checkErrors(m *raccoon.Infrastructure) (*raccoon.Infrastructure, error) {
	err := false
	if len(m.Infrastructure) == 0 {
		log.Error("No cluster were found on infrastructure file")
		err = true
	}

	if m.Name == "" {
		log.Errorf("infrastructure name can't be blank")
		err = true
	}

	for _, cluster := range m.Infrastructure {
		if len(cluster.Hosts) == 0 {
			log.Errorf("No hosts were found on cluster '%s' for commands '%s'",
				cluster.Name, cluster.TasksToExecute)
			err = true
		}

		if cluster.Name == "" {
			log.Errorf("hosts name can't be blank")
			err = true
		}

		if len(cluster.TasksToExecute) == 0 {
			log.Errorf("You haven't specified any task. Specify at least one as an string array on cluster '%s'", cluster.Name)
			err = true
		}

		for _, host := range cluster.Hosts {
			if host.Username == "" {
				log.Errorf("Host username is blank on host '%s'", host.IP)
				err = true
			}

			if host.Password == "" {
				log.Warnf("Host password is blank on host '%s'. If no password"+
					" is specified you must use an identity file or an interactive"+
					" authentication method", host.IP)
			}

			if host.IP == "" {
				log.Errorf("Host IP can't be blank on host '%s'", host.IP)
				err = true
			}
		}
	}
	if err {
		return &raccoon.Infrastructure{}, errors.New("Error found when parsing " +
			"infrastructure file")
	}

	return m, nil
}
