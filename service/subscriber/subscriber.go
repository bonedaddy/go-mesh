package subscriber

import (
	"context"

	logging "github.com/ipfs/go-log"
	inet "github.com/libp2p/go-libp2p-core/network"
	bc "github.com/upperwal/go-mesh/interface/blockchain"
	service "github.com/upperwal/go-mesh/interface/service"
)

var log = logging.Logger("svc-subscriber")

// Name and Version of this service.
const (
	NAME    = "subscriber"
	VERSION = 1
)

type TopicName string

type SubscriberService struct {
	service.ApplicationContainer

	blockchain bc.Blockchain

	// Pubsub related
	topicTracker map[TopicName]*TopicWrapper

	// Context
	ctx context.Context
}

func NewSubscriberService(b bc.Blockchain) *SubscriberService {
	ss := &SubscriberService{
		blockchain:   b,
		topicTracker: make(map[TopicName]*TopicWrapper),
	}

	ss.SetNameVersion(NAME, VERSION)

	return ss
}

func (subService *SubscriberService) Start(ctx context.Context) error {
	log.Info(subService.GetName(), "service started")

	subService.ctx = ctx

	log.Info(subService.GetName(), "sub")

	return nil
}

func (subService *SubscriberService) Stop() error {
	return nil
}

func (subService *SubscriberService) Run(stream inet.Stream) {

}

func (subService *SubscriberService) Get(name string) (chan interface{}, error) {
	return nil, nil
}

func (subService *SubscriberService) Set(name string, value interface{}) error {
	return nil
}