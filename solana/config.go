package solana

import "github.com/strangelove-ventures/noble-cctp-relayer/types"

var _ types.ChainConfig = (*Config)(nil)

type Config struct {
	RPC    string `yaml:"rpc"`
	WS     string `yaml:"ws"`
	Domain types.Domain

	MessageTransmitter   string                  `yaml:"message-transmitter"`
	TokenMessengerMinter string                  `yaml:"token-messenger-minter"`
	FiatToken            string                  `yaml:"fiat-token"`
	RemoteTokens         map[types.Domain]string `json:"remote-tokens"`

	StartBlock     uint64 `yaml:"start-block"`
	LookbackPeriod uint64 `yaml:"lookback-period"`

	BroadcastRetries       int `yaml:"broadcast-retries"`
	BroadcastRetryInterval int `yaml:"broadcast-retry-interval"`

	MinMintAmount uint64 `yaml:"min-mint-amount"`

	MetricsDenom    string `yaml:"metrics-denom"`
	MetricsExponent int    `yaml:"metrics-exponent"`

	MinterPrivateKey string `yaml:"minter-private-key"`
}

func (cfg Config) Chain(_ string) (types.Chain, error) {
	return NewSolana(cfg), nil
}
