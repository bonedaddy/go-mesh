package bootstrap

import (
	"time"

	addr "github.com/ipfs/go-ipfs-addr"
	inet "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

func (bs *BootstrapService) startBootstrapping() error {
	/* v1b := cid.V1Builder{
		Codec:  cid.Raw,
		MhType: mh.SHA2_256,
	}

	rendezvousPoint, _ := v1b.Sum([]byte(bs.rendezvousPoint)) */

	err := bs.connectToBootstrapPeers()
	if err != nil {
		return err
	}
	// TODO: Wait for connectToBootstrapPeers to complete before proceeding.
	err = bs.announceAndFind()
	if err != nil {
		return err
	}

	return nil
}

func (bs *BootstrapService) connectToBootstrapPeers() error {
	host := bs.GetHost()
	for _, peerAddr := range bs.bootstrapPeers {
		ipfsAddr, err := addr.ParseString(peerAddr)
		if err != nil {
			return err
		}

		peerInfo, err := pstore.InfoFromP2pAddr(ipfsAddr.Multiaddr())
		if err != nil {
			return err
		}

		if host.ID() == peerInfo.ID {
			continue
		}

		if err := host.Connect(bs.ctxLocal, *peerInfo); err != nil {
			log.Error("Could not connect to bootstrap peer: ", peerAddr, err)
			return err
		} else {
			log.Info("Connectivity with bootstrap peer success! Peer: ", peerAddr)
		}

		bs.s, err = host.NewStream(bs.ctxLocal, peerInfo.ID, bs.GetProtocol())
		if err != nil {
			return err
		}
		log.Info("Started new bootstrap stream")

		go bs.Run(bs.s)
		bs.WriteTemp()
	}

	return nil
}

func (bs *BootstrapService) announceAndFind() error {
	dht := bs.GetDHT()
	host := bs.GetHost()

	routingDiscovery := discovery.NewRoutingDiscovery(dht)
	discovery.Advertise(bs.ctxLocal, routingDiscovery, bs.rendezvousPoint)
	log.Debug("Successfully announced!")

	ticker := time.NewTicker(bs.d)

	for {
		log.Debug("Searching for other peers...")
		peerChan, err := routingDiscovery.FindPeers(bs.ctxLocal, bs.rendezvousPoint)
		if err != nil {
			return err
		}

		for p := range peerChan {
			if p.ID == host.ID() {
				continue
			}
			log.Debug("Found peer:", p)

			if host.Network().Connectedness(p.ID) == inet.Connected {
				log.Info("Already have a conn to", p.ID)
				continue
			}

			go func(pi peer.AddrInfo) {
				err := host.Connect(bs.ctxLocal, pi)

				if err != nil {
					log.Warning("Connection failed:", err)
					return
				}

				log.Info("Connected to:", pi)
			}(p)
		}

		select {
		case <-bs.ctxLocal.Done():
			ticker.Stop()
			return nil
		default:
		}

		<-ticker.C
	}
}
