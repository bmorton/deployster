package poller

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bmorton/deployster/schema"
)

const (
	// defaultTimeout is the amount of time to attempt checking for the new
	// version of the service to boot.  If checking exceeds this time, the
	// check will be aborted and any dependent actions will be cancelled.
	defaultTimeout time.Duration = 5 * time.Minute

	// defaultDelay is the amount of to wait between checks for the boot
	// completion of the new version.
	defaultDelay time.Duration = 1 * time.Second
)

type Poller struct {
	Deploy              *schema.Deploy
	Timeout             time.Duration
	Delay               time.Duration
	client              FleetClient
	stopChan            chan string
	successChan         chan *Event
	failureChan         chan *Event
	unresolvedChan      chan *Event
	successHandlers     []Handler
	unresolvedInstances map[string]*schema.ServiceInstance
}

func New(deploy *schema.Deploy, client FleetClient) *Poller {
	toBeResolved := make(map[string]*schema.ServiceInstance)
	for i := 1; i <= deploy.InstanceCount; i++ {
		instance := deploy.ServiceInstance(strconv.Itoa(i))
		toBeResolved[instance.FleetUnitName()] = instance
	}
	return &Poller{
		Deploy:              deploy,
		Timeout:             defaultTimeout,
		Delay:               defaultDelay,
		stopChan:            make(chan string, 1),
		successChan:         make(chan *Event, deploy.InstanceCount),
		failureChan:         make(chan *Event, deploy.InstanceCount),
		unresolvedChan:      make(chan *Event, deploy.InstanceCount),
		client:              client,
		unresolvedInstances: toBeResolved,
	}
}

func (p *Poller) Watch() {
	timeout := time.After(p.Timeout)

	for {
		if len(p.unresolvedInstances) == 0 {
			return
		}
		pollStates := time.After(p.Delay)

		select {
		case <-pollStates:
			p.pollStates()
		case event := <-p.successChan:
			log.Printf("%s is running.\n", event.ServiceInstance.FleetUnitName())
			p.runSuccessHandlers(event)
			delete(p.unresolvedInstances, event.ServiceInstance.FleetUnitName())
		case event := <-p.failureChan:
			log.Printf("%s failed to launch.\n", event.ServiceInstance.FleetUnitName())
			delete(p.unresolvedInstances, event.ServiceInstance.FleetUnitName())
		case event := <-p.unresolvedChan:
			log.Printf("%s is not yet resolved (state: %s).\n", event.ServiceInstance.FleetUnitName(), event.SystemdSubState)
		case <-timeout:
			p.stopChan <- fmt.Sprintf("Timed out polling state of %s:%s after %s.\n", p.Deploy.ServiceName, p.Deploy.Version, p.Timeout)
		case msg := <-p.stopChan:
			log.Println(msg)
			return
		}
	}

	return
}

func (p *Poller) AddSuccessHandler(newHandler Handler) {
	p.successHandlers = append(p.successHandlers, newHandler)
}

func (p *Poller) runSuccessHandlers(event *Event) {
	for _, h := range p.successHandlers {
		h.Handle(event)
	}
	return
}

func (p *Poller) pollStates() {
	log.Printf("Checking state(s) of %s:%s...\n", p.Deploy.ServiceName, p.Deploy.Version)
	events, err := p.fetchStates()
	if err != nil {
		log.Println(err)
		return
	}

	for _, event := range events {
		switch event.SystemdSubState {
		case "running":
			p.successChan <- event
		case "failed":
			p.failureChan <- event
		default:
			p.unresolvedChan <- event
		}
	}

	return
}

func (p *Poller) fetchStates() ([]*Event, error) {
	var events []*Event
	states, err := p.client.UnitStates()
	if err != nil {
		return events, err
	}

	for _, state := range states {
		if instance, ok := p.unresolvedInstances[state.Name]; ok {
			events = append(events, NewEvent(instance, state))
		}
	}

	return events, nil
}
