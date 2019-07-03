package lnd

import (
	"context"

	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-bitcoin/lnrpc"
)

func (client *Client) SetChannelHandler(channelHandler bitcoin.ChannelHandler) {
	client.channelHandler = channelHandler
}

func (client *Client) handleChannelEvent(event *lnrpc.ChannelEventUpdate) {
	if client.channelHandler == nil {
		return
	}

	ch := &bitcoin.Channel{}

	switch event.GetType() {
	case lnrpc.ChannelEventUpdate_OPEN_CHANNEL:
		oc := event.GetOpenChannel()
		ch.Balance = oc.GetLocalBalance()
		ch.Capacity = oc.GetCapacity()
		ch.Online = oc.GetActive()
		ch.ShortChannelID = BOLT7Uint64ToString(oc.GetChanId())
		ch.RemotePublicKey = oc.GetRemotePubkey()
		ch.State = bitcoin.ChannelStateOpen

		client.channelHandler.ChannelOpened(ch)
	}

}

func (client *Client) subscribeChannels() {
	s, err := client.rpc.SubscribeChannelEvents(context.Background(), &lnrpc.ChannelEventSubscription{})
	if err != nil {
		panic(err)
	}
	for {
		update, err := s.Recv()
		if err != nil {
			return
		}

		client.handleChannelEvent(update)
	}
}
