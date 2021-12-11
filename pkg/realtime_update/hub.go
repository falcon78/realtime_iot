package realtime_update

import "time"

type Payload struct {
	ChannelOne   float64   `json:"channelOne"`
	ChannelTwo   float64   `json:"channelTwo"`
	ChannelThree float64   `json:"channelThree"`
	ChannelFour  float64   `json:"channelFour"`
	Timestamp    time.Time `json:"timestamp"`
}

type Message struct {
	AccessKey string
	Payload   *Payload
}

type SubscriptionType struct {
	ChannelName string
	Client      *Client
}

type Hub struct {
	Subscriptions map[string]map[*Client]bool
	Broadcast     chan *Message
	Register      chan *SubscriptionType
	Unregister    chan *SubscriptionType
}

func New() *Hub {
	return &Hub{
		Broadcast:     make(chan *Message),
		Register:      make(chan *SubscriptionType),
		Unregister:    make(chan *SubscriptionType),
		Subscriptions: make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Listen() {
	for {
		select {
		case subscription := <-h.Register:
			channelSubscription, ok := h.Subscriptions[subscription.ChannelName]
			if !ok {
				h.Subscriptions[subscription.ChannelName] = make(map[*Client]bool)
				channelSubscription = h.Subscriptions[subscription.ChannelName]
			}
			channelSubscription[subscription.Client] = true
		case subscription := <-h.Unregister:
			if _, ok := h.Subscriptions[subscription.ChannelName]; ok {
				if _, ok := h.Subscriptions[subscription.ChannelName][subscription.Client]; ok {
					delete(h.Subscriptions[subscription.ChannelName], subscription.Client)
					close(subscription.Client.Send)
				}
			}
		case message := <-h.Broadcast:
			for client := range h.Subscriptions[message.AccessKey] {
				select {
				case client.Send <- message.Payload:
				default:
					close(client.Send)
					delete(h.Subscriptions[message.AccessKey], client)
				}
			}
		}
	}
}
