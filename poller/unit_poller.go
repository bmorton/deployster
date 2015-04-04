package poller

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bmorton/deployster/schema"
	fleet "github.com/coreos/fleet/schema"
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

type Handler interface {
	Handle(*Event)
}

type UnitPoller struct {
	ServiceInstance *schema.ServiceInstance
	Timeout         time.Duration
	Delay           time.Duration
	client          FleetClient
	stopChan        chan string
	successChan     chan *Event
	failureChan     chan *Event
	unresolvedChan  chan *Event
	successHandlers []Handler
}

func New(serviceInstance *schema.ServiceInstance, client FleetClient) *UnitPoller {
	return &UnitPoller{
		ServiceInstance: serviceInstance,
		Timeout:         defaultTimeout,
		Delay:           defaultDelay,
		stopChan:        make(chan string, 1),
		successChan:     make(chan *Event, 1),
		failureChan:     make(chan *Event, 1),
		unresolvedChan:  make(chan *Event, 1),
		client:          client,
	}
}

func (p *UnitPoller) Watch() {
	timeout := time.After(p.Timeout)

	for {
		pollStatus := time.After(p.Delay)

		select {
		case <-pollStatus:
			p.handleStatus()
		case unit := <-p.successChan:
			log.Printf("%s is running.\n", p.ServiceInstance.FleetUnitName())
			p.runSuccessHandlers(unit)
			return
		case <-p.failureChan:
			log.Printf("%s failed to launch, bailing out.\n", p.ServiceInstance.FleetUnitName())
			return
		case state := <-p.unresolvedChan:
			log.Printf("%s is not yet resolved (state: %s).  Trying again in %s.\n", p.ServiceInstance.FleetUnitName(), state.SystemdSubState, p.Delay)
		case <-timeout:
			p.stopChan <- fmt.Sprintf("Timed out polling state of %s after %s.\n", p.ServiceInstance.FleetUnitName(), p.Timeout)
		case msg := <-p.stopChan:
			log.Println(msg)
			return
		}
	}

	return
}

func (p *UnitPoller) AddSuccessHandler(newHandler Handler) {
	p.successHandlers = append(p.successHandlers, newHandler)
}

func (p *UnitPoller) runSuccessHandlers(event *Event) {
	for _, h := range p.successHandlers {
		h.Handle(event)
	}
	return
}

func (p *UnitPoller) handleStatus() {
	log.Printf("Checking if %s has finished launching...\n", p.ServiceInstance.FleetUnitName())
	event, err := p.fetchStatus()
	if err != nil {
		log.Println(err)
		p.unresolvedChan <- event
		return
	}

	switch event.SystemdSubState {
	case "running":
		p.successChan <- event
	case "failed":
		p.failureChan <- event
	default:
		p.unresolvedChan <- event
	}

	return
}

func (p *UnitPoller) fetchStatus() (*Event, error) {
	states, err := p.client.UnitStates()
	if err != nil {
		return NewEvent(p.ServiceInstance, &fleet.UnitState{}), err
	}

	for _, state := range states {
		if state.Name == p.ServiceInstance.FleetUnitName() {
			return NewEvent(p.ServiceInstance, state), nil
		}
	}

	return NewEvent(p.ServiceInstance, &fleet.UnitState{}), errors.New("Unit state couldn't be determined")
}
