package dispatcher

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon/job"
	"github.com/thehivecorporation/raccoon/constants"
)

var wg sync.WaitGroup

//Dispatch receives a Job struct and is the one in charge of iterating for each
//Node within the job cluster of nodes and open a new routine for each of them
//with the recipe to execute.
func Dispatch(js *[]job.Job) {
	for _, j := range *js {
		log.WithFields(log.Fields{
			constants.HOST_NAME:j.Cluster.Name,
			constants.GROUP_NAME:j.GroupName,
		}).Info(constants.ARROW_LENGTH + constants.HOST_LAUNCH_MESSAGE)

		for _, node := range j.Cluster.Nodes {
			wg.Add(len(j.Zbook.Instructions))
			go executeRecipeOnNode(*j.Zbook, node)
		}
	}

	wg.Wait()
}

//executeRecipeOnNode will take every instruction of the recipe and execute it
//in order on each node. Each instruction waits until previous one is finished.
func executeRecipeOnNode(r job.Zbook, n connection.Node) {
	for _, instruction := range r.Instructions {
		instruction.Execute(n)
		wg.Done()
	}
}
