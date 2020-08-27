package standriver

import (
	"github.com/RezaOptic/go-utils/logger"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

// SubscribeFunc takes a STAN connection and subscribes anything it wants to it.
type SubscriberFunc func(sc stan.Conn)

type ConnectionInterface interface {
	Connect() error
	Reconnect() error
	AutoReconnect(interval time.Duration, stop chan struct{})
	Close() error
}

// Connection holds NATS and NATS-Streaming(STAN) connections.
// It's fields are meant to be read-only, so you should not write to them. To create a new Connection, you must use the
// standriver.New(...) function. Instantiating a Connection literal or using go's new() function may lead to undefined behavior.
type Connection struct {
	url                  string
	clusterID            string
	clientID             string
	HealthCheckerCreated bool
	Reconnected          bool
	subscribersSet       map[*SubscriberFunc]bool

	nats *nats.Conn
	stan stan.Conn
}

// New creates a new Connection
func New(URLs []string, ClusterID string, ClientID string) *Connection {
	x := &Connection{
		url:            strings.Join(URLs, ","),
		clusterID:      ClusterID,
		clientID:       ClientID,
		subscribersSet: make(map[*SubscriberFunc]bool),
	}

	return x
}

// Connect connects to NATS and STAN, initializing their corresponding variables
func (n *Connection) Connect() (err error) {
	defer func() {
		_ = logger.ZSLogger.Sync()
	}()

	logger.ZSLogger.Infof("connecting to NATS host/cluster: %v", n.url)

	n.nats, err = nats.Connect(n.url,
		nats.Timeout(10*time.Second),
		nats.PingInterval(2*time.Second),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(1*time.Second),
		nats.ReconnectBufSize(8*1024*1024),
	)

	//)
	if err != nil {
		logger.ZSLogger.Errorf("can not connect to NATS server: %v", err)
		return err
	}
	logger.ZSLogger.Infof("connected to NATS")

	n.stan, err = stan.Connect(n.clusterID, n.clientID, stan.NatsURL(n.url), stan.Pings(10, 5),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			logger.ZSLogger.Errorf("Connection lost, reason: %v", err)
		}))
	if err != nil {
		logger.ZSLogger.Errorf(`can not connect to STAN server with error : %s`, err)
		return err
	}
	logger.ZSLogger.Infof("connected to STAN")

	//// call all SubscriberFunctions on connect/reconnect
	for sf := range n.subscribersSet {
		(*sf)(n.stan)
	}
	if !n.HealthCheckerCreated {
		go func() {
			for {
				if (n.nats == nil || n.nats.Status() != nats.CONNECTED || n.stan == nil) && ! n.Reconnected {
					errReconnect := n.Reconnect()
					if errReconnect != nil {
						n.Reconnected = false
					}
				}
				time.Sleep(time.Second * 5)
			}
		}()
	}
	n.HealthCheckerCreated = true
	n.Reconnected = false
	return nil
}

// Reconnect closes STAN connection and connects again
func (n *Connection) Reconnect() (err error) {
	n.Reconnected = true
	defer func() {
		_ = logger.ZSLogger.Sync()
	}()

	logger.ZSLogger.Infof("reconnecting to STAN")
	if n.stan != nil {
		logger.ZSLogger.Infof("CLOSING to STAN")
		_ = n.stan.Close()
	}

	return n.Connect()
}

// Close closes NATS and STAN connections
func (n *Connection) Close() (err error) {
	defer func() {
		_ = logger.ZSLogger.Sync()
	}()

	err = n.stan.Close()
	if err != nil {
		logger.ZSLogger.Errorf("can not close NATS Streaming connection: %v", err)
	}
	n.nats.Close()

	return err
}

// Publish publishes the message with subject
func (n *Connection) Publish(subject string, message []byte) (err error) {
	// TODO handle n.stan==nil case
	return n.stan.Publish(subject, message)
}

// Subscribe takes a SubscriberFunc sf and invokes it, and also makes sure it is called every time STAN is reconnected.
// This means subscriptions persist through STAN reconnects.
// SubscriberFunc is a function which performs your subscriptions to a stan.Conn object.
func (n *Connection) Subscribe(sf SubscriberFunc) {
	_, ok := n.subscribersSet[&sf]
	if ok {
		// already subscribed
		return
	}

	if n.stan != nil {
		sf(n.stan)
	}

	n.subscribersSet[&sf] = true
}

func (n *Connection) reconnectedHandler(nc *nats.Conn) {
	logger.ZSLogger.Infof("reconnected to NATS")
}
