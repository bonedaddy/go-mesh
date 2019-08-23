package mesh

import (
	"context"
	"fmt"

	application "github.com/upperwal/go-mesh/application"
	config "github.com/upperwal/go-mesh/config"
	driver "github.com/upperwal/go-mesh/driver/eth"
	ra "github.com/upperwal/go-mesh/interface/ra"
	fpubsub "github.com/upperwal/go-mesh/pubsub"
	bs "github.com/upperwal/go-mesh/service/bootstrap"
	pubservice "github.com/upperwal/go-mesh/service/publisher"
	subservice "github.com/upperwal/go-mesh/service/subscriber"
)

// Mesh contains variables need for client to work.
type Mesh struct {
	Cfg          *config.Config
	app          *application.Application
	remoteAccess ra.RemoteAccess
	bootService  *bs.BootstrapService

	// Data Services
	SubService *subservice.SubscriberService
	PubService *pubservice.PublisherService
}

// NewMesh creates a new mesh client.
func NewMesh(ctx context.Context, opt ...config.Option) (*Mesh, error) {

	cfg, err := config.NewConfig(opt...)
	if err != nil {
		return nil, err
	}

	if err = checkConfig(cfg); err != nil {
		return nil, err
	}

	app, err := application.NewApplication(
		ctx,
		cfg.AccountPrvKey,
		cfg.AccountCert,
		cfg.Host,
		cfg.Port,
	)
	if err != nil {
		return nil, err
	}

	ethDriver, err := driver.NewEthDriver(cfg.RemoteAccessURI, cfg.RemoteAccessPrivateKey)
	if err != nil {
		return nil, err
	}

	bservice := bs.NewBootstrapService(false, cfg.BootstrapRendezvous, cfg.BootstrapNodes, cfg.BootstrapRefreshRate)
	app.InjectService(bservice)

	f := fpubsub.NewFilter(ethDriver)
	app.SetGossipPeerFilter(f)

	return &Mesh{
		Cfg:          cfg,
		app:          app,
		remoteAccess: ethDriver,
		bootService:  bservice,
	}, nil
}

// Start the mesh client and all its services.
func (m *Mesh) Start() error {
	fmt.Println("Starting... mesh")
	err := m.app.Start()
	if err != nil {
		return err
	}
	return nil
}

// AddSubscriber adds subscriber service to the mesh
func (m *Mesh) AddSubscriber() {
	subService := subservice.NewSubscriberService(m.remoteAccess)
	m.app.InjectService(subService)

	m.SubService = subService
}

// AddPublisher adds publisher service to the mesh
func (m *Mesh) AddPublisher() {
	pubService := pubservice.NewPublisherService(m.remoteAccess)
	m.app.InjectService(pubService)

	m.PubService = pubService
}

// Wait blocks 'this' thread.
func (m *Mesh) Wait() {
	m.app.Wait()
}
