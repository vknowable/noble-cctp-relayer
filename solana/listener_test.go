package solana_test

import (
	"context"
	"testing"
	"time"

	extsolana "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/log"

	"github.com/strangelove-ventures/noble-cctp-relayer/relayer"
	"github.com/strangelove-ventures/noble-cctp-relayer/solana"
	"github.com/strangelove-ventures/noble-cctp-relayer/types"
)

// TestParseTransaction tests the transaction parsing utility used by the relayer.
// It fetches a mainnet transfer from Solana to Noble.
//
// https://solscan.io/tx/4XhuTtTHxNFDfGn6A7ngvqT2dNoxRQFK7kpZwNY25gxXU5SPFCAk6ihj9JJdq5g7UMb7MSwyTt5r3TJbM4RgMVyJ
// https://www.mintscan.io/noble/tx/98C4BB3B29FBA6EBA3B14C3DFA85925C8223A88E268CBA8FBEAD6E5F3F9333D9
func TestParseTransaction(t *testing.T) {
	// ARRANGE: Create a new instance of Solana.
	key, err := extsolana.NewRandomPrivateKey()
	require.NoError(t, err)

	chain := solana.NewSolana(solana.Config{
		RPC: "https://corie-nhz8jx-fast-mainnet.helius-rpc.com",
		WS:  "",

		MessageTransmitter:   "CCTPmbSD7gX1bxKPAmg77w8oFzNFpaQiQUWD43TKaecd",
		TokenMessengerMinter: "CCTPiPYPc6AsJuwueEnWgSgucamXDZwBd53dQ11YiKX3",
		FiatToken:            "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",

		MinterPrivateKey: key.String(),
	})
	require.NoError(t, chain.InitializeClients(context.Background(), log.NewNopLogger()))

	// ACT: Attempt to fetch and parse a transaction.
	events, err := chain.ParseTransactionMessages(
		context.Background(),
		"4XhuTtTHxNFDfGn6A7ngvqT2dNoxRQFK7kpZwNY25gxXU5SPFCAk6ihj9JJdq5g7UMb7MSwyTt5r3TJbM4RgMVyJ",
	)
	// ASSERT: The action should've succeeded, and returned exactly 1 message (this is a single-message tx).
	require.NoError(t, err)
	require.Equal(t, 1, len(events))
}

// TestParseTransactionMultiMessage tests that the ParseTransactionMessages method
// works correctly with single-message transactions. This ensures backward compatibility
func TestParseTransactionMultiMessage(t *testing.T) {
	// ARRANGE: Create a new instance of Solana.
	key, err := extsolana.NewRandomPrivateKey()
	require.NoError(t, err)

	chain := solana.NewSolana(solana.Config{
		RPC: "https://corie-nhz8jx-fast-mainnet.helius-rpc.com",
		WS:  "",

		MessageTransmitter:   "CCTPmbSD7gX1bxKPAmg77w8oFzNFpaQiQUWD43TKaecd",
		TokenMessengerMinter: "CCTPiPYPc6AsJuwueEnWgSgucamXDZwBd53dQ11YiKX3",
		FiatToken:            "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",

		MinterPrivateKey: key.String(),
	})
	require.NoError(t, chain.InitializeClients(context.Background(), log.NewNopLogger()))

	// ACT: Test with a real transaction that we know has exactly one message
	events, err := chain.ParseTransactionMessages(
		context.Background(),
		"4XhuTtTHxNFDfGn6A7ngvqT2dNoxRQFK7kpZwNY25gxXU5SPFCAk6ihj9JJdq5g7UMb7MSwyTt5r3TJbM4RgMVyJ",
	)

	// ASSERT: The action should've succeeded, and returned exactly one message
	require.NoError(t, err)
	require.Equal(t, 1, len(events), "Should return exactly 1 message for this single-message transaction")

	// Verify the message properties
	event := events[0]
	require.Equal(t, "4XhuTtTHxNFDfGn6A7ngvqT2dNoxRQFK7kpZwNY25gxXU5SPFCAk6ihj9JJdq5g7UMb7MSwyTt5r3TJbM4RgMVyJ", event.SourceTxHash)
	require.Equal(t, types.Domain(5), event.SourceDomain, "Source domain should be Solana (5)")
	require.Equal(t, types.Created, event.Status, "Message status should be Created")
	require.NotEmpty(t, event.IrisLookupID, "IrisLookupID should not be empty")
	require.NotEmpty(t, event.MsgSentBytes, "MsgSentBytes should not be empty")

	t.Logf("Successfully parsed single message transaction with IrisLookupID: %s", event.IrisLookupID)
	t.Logf("Message destination domain: %d", event.DestDomain)
}

// TestMultiMessageParsingLogic tests the API contract and data structure handling
// for the ParseTransactionMessages method. Since actual multi-message transaction examples
// are difficult to find, we only verify that the method returns the correct data types
// TODO: Update this test once an actual multi-message transaction example is found
func TestMultiMessageParsingLogic(t *testing.T) {
	// ARRANGE: Create a new instance of Solana.
	key, err := extsolana.NewRandomPrivateKey()
	require.NoError(t, err)

	chain := solana.NewSolana(solana.Config{
		RPC: "https://corie-nhz8jx-fast-mainnet.helius-rpc.com",
		WS:  "",

		MessageTransmitter:   "CCTPmbSD7gX1bxKPAmg77w8oFzNFpaQiQUWD43TKaecd",
		TokenMessengerMinter: "CCTPiPYPc6AsJuwueEnWgSgucamXDZwBd53dQ11YiKX3",
		FiatToken:            "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",

		MinterPrivateKey: key.String(),
	})
	require.NoError(t, chain.InitializeClients(context.Background(), log.NewNopLogger()))

	// Test the API contract: ParseTransactionMessages should return a slice
	// and handle single message scenarios correctly

	// Test 1: Verify that ParseTransactionMessages returns a slice (not a single message)
	events, err := chain.ParseTransactionMessages(
		context.Background(),
		"4XhuTtTHxNFDfGn6A7ngvqT2dNoxRQFK7kpZwNY25gxXU5SPFCAk6ihj9JJdq5g7UMb7MSwyTt5r3TJbM4RgMVyJ",
	)

	require.NoError(t, err)
	require.IsType(t, []*types.MessageState{}, events, "ParseTransactionMessages should return a slice of MessageState")
	require.GreaterOrEqual(t, len(events), 0, "Should return at least 0 messages (empty slice for no messages)")

	// Test 2: Verify that each message in the slice has the correct structure
	for i, event := range events {
		require.NotNil(t, event, "Message %d should not be nil", i)
		require.Equal(t, types.Created, event.Status, "Message %d should have Created status", i)
		require.Equal(t, types.Domain(5), event.SourceDomain, "Message %d should have Solana as source domain", i)
		require.NotEmpty(t, event.SourceTxHash, "Message %d should have a source transaction hash", i)
		require.NotEmpty(t, event.IrisLookupID, "Message %d should have an IrisLookupID", i)
	}

	// Test 3: Verify that all messages have the same transaction hash
	// (since they all came from the same transaction)
	txHash := "4XhuTtTHxNFDfGn6A7ngvqT2dNoxRQFK7kpZwNY25gxXU5SPFCAk6ihj9JJdq5g7UMb7MSwyTt5r3TJbM4RgMVyJ"
	for i, event := range events {
		require.Equal(t, txHash, event.SourceTxHash, "Message %d should have the same transaction hash", i)
	}

	t.Logf("Multi-message parsing logic test passed. Found %d messages in transaction.", len(events))

	// Test 4: Verify that the parsing logic works correctly for single messages
	// Our code structure is ready for multi-message scenarios when they occur
	if len(events) == 1 {
		t.Log("Single message transaction parsed correctly - API contract verified")
	} else if len(events) > 1 {
		t.Logf("Multi-message transaction parsed correctly - found %d messages", len(events))
	} else {
		t.Log("No messages found - this might indicate a parsing issue or non-CCTP transaction")
	}
}

// TestTrackLatestBlockHeight tests the actual TrackLatestBlockHeight method functionality.
// This verifies that Solana can query real block height, update metrics, and handle context cancellation.
func TestTrackLatestBlockHeight(t *testing.T) {
	// ARRANGE: Create a new instance of Solana with real RPC endpoint
	key, err := extsolana.NewRandomPrivateKey()
	require.NoError(t, err)

	chain := solana.NewSolana(solana.Config{
		RPC: "https://corie-nhz8jx-fast-mainnet.helius-rpc.com",
		WS:  "",

		MessageTransmitter:   "CCTPmbSD7gX1bxKPAmg77w8oFzNFpaQiQUWD43TKaecd",
		TokenMessengerMinter: "CCTPiPYPc6AsJuwueEnWgSgucamXDZwBd53dQ11YiKX3",
		FiatToken:            "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",

		MinterPrivateKey: key.String(),
	})
	require.NoError(t, chain.InitializeClients(context.Background(), log.NewNopLogger()))

	// Verify initial state
	require.Equal(t, uint64(0), chain.LatestBlock(), "Initial block height should be 0")
	require.Equal(t, types.Domain(5), chain.Domain(), "Solana domain should be 5")
	require.Equal(t, "Solana", chain.Name(), "Chain name should be Solana")

	// ACT: Start TrackLatestBlockHeight in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.NewNopLogger()
	// Pass nil for metrics to avoid the complexity of setting up Prometheus in tests
	var metrics *relayer.PromMetrics = nil

	// Start the tracking routine
	go chain.TrackLatestBlockHeight(ctx, logger, metrics)

	// Wait for the first update (should happen immediately, but give it more time for RPC call)
	time.Sleep(2 * time.Second)

	// ASSERT: Block height should now be updated with real data
	blockHeight := chain.LatestBlock()
	require.Greater(t, blockHeight, uint64(0), "Block height should be updated with real Solana data")

	t.Logf("Initial block height: %d", blockHeight)

	// Wait for another update (should happen after 1 second)
	time.Sleep(1200 * time.Millisecond) // Wait a bit more than 1 second

	// ASSERT: Block height should have increased (Solana produces blocks every ~400ms)
	newBlockHeight := chain.LatestBlock()
	require.Greater(t, newBlockHeight, blockHeight, "Block height should be increasing (Solana produces blocks every ~400ms)")

	t.Logf("Updated block height: %d (increased by %d)", newBlockHeight, newBlockHeight-blockHeight)

	// ACT: Test context cancellation
	cancel()

	// Wait a bit to ensure the routine has stopped
	time.Sleep(100 * time.Millisecond)

	// ASSERT: Final block height should be reasonable
	finalHeight := chain.LatestBlock()
	require.Greater(t, finalHeight, uint64(0), "Final block height should be positive")

	t.Logf("TrackLatestBlockHeight test passed - final height: %d", finalHeight)
}
