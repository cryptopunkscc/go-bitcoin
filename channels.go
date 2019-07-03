package bitcoin

type ChannelState int

const (
	ChannelStateUnknown      ChannelState = iota // Channel state is unknown
	ChannelStateOpening                          // Waiting for openening transaction confirmation
	ChannelStateOpen                             // Channel is open
	ChannelStateClosing                          // Channel is being cooperatively closed
	ChannelStateForceClosing                     // Channel is being forcefully closed
	ChannelStateCloseWait                        // Waiting for closing transaction confirmation
)

type Channel struct {
	State           ChannelState // State of the channel
	RemotePublicKey string       // Public key of the remote node
	FundingTxID     string       // TxID of the funding transaction
	ShortChannelID  string       // BOLT7 channel ID
	Capacity        int64        // Total capacity of the channel
	Balance         int64        // Local balance of the channel
	Online          bool         // Are we currently connected to the peer
}

type ChannelHandler interface {
	ChannelOpened(*Channel)
	ChannelClosed(*Channel)
}
