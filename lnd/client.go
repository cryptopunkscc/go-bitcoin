package lnd

import (
	"context"

	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-bitcoin/lnrpc"
	"google.golang.org/grpc"
)

type Client struct {
	rpc            lnrpc.LightningClient
	channelHandler bitcoin.ChannelHandler
	invoiceHandler bitcoin.InvoiceHandler
}

func Connect(cfg *Config) (*Client, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(cfg.TLSCredentials()),
		grpc.WithPerRPCCredentials(cfg.Macaroon()),
	}

	conn, err := grpc.Dial(cfg.url(), opts...)

	if err != nil {
		return nil, err
	}

	c := &Client{
		rpc: lnrpc.NewLightningClient(conn),
	}

	go c.subscribeChannels()
	go c.subscribeInvoices()

	return c, nil
}

func (client *Client) Version() string {
	res, err := client.rpc.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})

	if err != nil {
		panic(err)
	}

	return "lnd " + res.Version
}

func (client *Client) Balance() int {
	res, _ := client.rpc.ChannelBalance(context.Background(), &lnrpc.ChannelBalanceRequest{})

	return int(res.Balance)
}
