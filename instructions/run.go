package instructions

import (
	log "github.com/Sirupsen/logrus"

	"github.com/thehivecorporation/raccoon"
)

//RUN is a instruction that in the recipe file correspond to the CMD instruction.
//It will execute the "Command" on every machine. Ideally, every command must
//be bash
type RUN struct {
	//The name that identifies this struct ("RUN" in this case)
	Name string

	//Description of the instruction that must be set by the user
	Description string

	//Bash instruction to execute
	Instruction string
}

func (r *RUN) GetCommandName() string {
	return "RUN"
}

//Execute is the implementation of the Instruction interface for a RUN instruction
func (r *RUN) Execute(h raccoon.Host) {
	session, err := h.GetSession()

	if err != nil {
		logError(err, r, &h)

		session.Close()

		return
	}
	defer session.Close()

	r.LogCommand(&h)

	if err = session.Run(r.Instruction); err != nil {
		logError(err, r, &h)
	}
}
func (r *RUN) LogCommand(h *raccoon.Host) {
	log.WithFields(log.Fields{
		"Instruction": r.Name,
		"Node":        h.IP,
		"package":     packageName,
	}).Info(r.Description)
}
