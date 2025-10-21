package cmd_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/strangelove-ventures/noble-cctp-relayer/cmd"
	"github.com/strangelove-ventures/noble-cctp-relayer/ethereum"
	"github.com/strangelove-ventures/noble-cctp-relayer/noble"
	"github.com/strangelove-ventures/noble-cctp-relayer/solana"
	"github.com/strangelove-ventures/noble-cctp-relayer/types"
)

func TestConfig(t *testing.T) {
	file, err := cmd.ParseConfig("../config/sample-config.yaml")
	require.NoError(t, err, "Error parsing config")

	// assert noble chainConfig correctly parsed
	var nobleType any = file.Chains["noble"]
	_, ok := nobleType.(*noble.ChainConfig)
	require.True(t, ok)

	// assert ethereum chainConfig correctly parsed
	var ethType any = file.Chains["ethereum"]
	_, ok = ethType.(*ethereum.ChainConfig)
	require.True(t, ok)

	// assert solana chainConfig correctly parsed
	var solType any = file.Chains["solana"]
	_, ok = solType.(*solana.Config)
	require.True(t, ok)
}

func TestSolanaConfigParsing(t *testing.T) {
	// Test that Solana config parsing works correctly
	// This test verifies that Solana chains are properly identified and parsed

	// Create a test config file content
	testConfigYAML := `
chains:
  solana:
    rpc: "https://corie-nhz8jx-fast-mainnet.helius-rpc.com"
    ws: "wss://corie-nhz8jx-fast-mainnet.helius-rpc.com"
    domain: 5
    message-transmitter: "CCTPmbSD7gX1bxKPAmg77w8oFzNFpaQiQUWD43TKaecd"
    token-messenger-minter: "CCTPiPYPc6AsJuwueEnWgSgucamXDZwBd53dQ11YiKX3"
    fiat-token: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
    remote-tokens:
      0: "0x487039debedbf32d260137b0a6f66b90962bec777250910d253781de326a716d"
    start-block: 0
    lookback-period: 5
    broadcast-retries: 5
    broadcast-retry-interval: 10
    min-mint-amount: 1000000
    metrics-denom: "SOL"
    metrics-exponent: 9
    minter-private-key: "test-key"

enabled-routes:
  5: [4] # solana -> noble

circle:
  attestation-base-url: "https://iris-api-sandbox.circle.com/attestations/"
  fetch-retries: 30
  fetch-retry-interval: 3

processor-worker-count: 1
`

	// Write test config to temporary file
	tmpFile := "/tmp/test-solana-config.yaml"
	err := os.WriteFile(tmpFile, []byte(testConfigYAML), 0644)
	require.NoError(t, err, "Failed to write test config file")
	defer os.Remove(tmpFile)

	// Parse the config
	config, err := cmd.ParseConfig(tmpFile)
	require.NoError(t, err, "Error parsing Solana config")

	// Verify Solana chain was parsed correctly
	solanaChain, exists := config.Chains["solana"]
	require.True(t, exists, "Solana chain should exist in parsed config")

	// Verify it's the correct type
	solanaConfig, ok := solanaChain.(*solana.Config)
	require.True(t, ok, "Solana chain should be parsed as solana.Config")

	// Verify some key fields
	require.Equal(t, "https://corie-nhz8jx-fast-mainnet.helius-rpc.com", solanaConfig.RPC)
	require.Equal(t, types.Domain(5), solanaConfig.Domain)
	require.Equal(t, "CCTPmbSD7gX1bxKPAmg77w8oFzNFpaQiQUWD43TKaecd", solanaConfig.MessageTransmitter)
}

func TestBlockQueueChannelSize(t *testing.T) {
	file, err := cmd.ParseConfig("../config/sample-config.yaml")
	require.NoError(t, err, "Error parsing config")

	var nobleCfg any = file.Chains["noble"]
	n, ok := nobleCfg.(*noble.ChainConfig)
	require.True(t, ok)

	// block-queue-channel-size is set to 1000000 in sample-config
	expected := uint64(1000000)

	require.Equal(t, expected, n.BlockQueueChannelSize)
}
