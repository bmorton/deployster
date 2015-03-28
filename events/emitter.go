package events

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

const defaultTTL = 60 * 60 * 24

type Emitter struct {
	client *etcd.Client
}

type Event interface {
	EventType() string
}

func NewEmitter() *Emitter {
	return &Emitter{client: etcd.NewClient([]string{os.Getenv("ETCDCTL_PEERS")})}
}

func (e *Emitter) Emit(event Event) {
	b, _ := json.Marshal(event)

	timestamp := time.Now().Unix()
	_, err := e.client.Set(fmt.Sprintf("/deployster/events/%s/%d", event.EventType(), timestamp), string(b), defaultTTL)
	if err != nil {
		log.Println(err)
	}
}
