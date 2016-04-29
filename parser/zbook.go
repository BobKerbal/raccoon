package parser

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon/instructions"
	"github.com/thehivecorporation/raccoon/job"
	"github.com/thehivecorporation/raccoon/constants"
)

type book struct {
	Title           string              `json:"book_title"`
	Maintainer      string              `json:"maintainer"`
	RawInstructions []map[string]string `json:"instructions"`
}

//ReadZbookFile will take a filepath as parameter and return a Job
func ReadZbookFile(f string) (job.Zbook, error) {
	log.WithFields(log.Fields{
		constants.INSTRUCTIONS_NAME: f,
	}).Info(constants.ARROW_LENGTH + "Reading " + constants.INSTRUCTIONS_NAME +
		" file")

	var z book

	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return job.Zbook{}, err
	}

	err = json.Unmarshal(dat, &z)
	if err != nil {
		return job.Zbook{}, err
	}

	parsedInstructions := make([]instructions.InstructionExecutor, 0)

	for _, i := range z.RawInstructions {
		switch i["name"] {
		case "RUN":
			run := instructions.RUN{
				Name:        "RUN",
				Description: i["description"],
				Instruction: i["instruction"],
			}
			parsedInstructions = append(parsedInstructions, &run)
		case "ADD":
			add := instructions.ADD{
				SourcePath:  i["sourcePath"],
				DestPath:    i["destPath"],
				Description: i["description"],
				Name:        "ADD",
			}
			parsedInstructions = append(parsedInstructions, &add)
		}
	}

	return job.Zbook{Instructions: parsedInstructions}, nil
}
