package testutil

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/log"

	extsolana "github.com/gagliardetto/solana-go"
	"github.com/strangelove-ventures/noble-cctp-relayer/cmd"
	"github.com/strangelove-ventures/noble-cctp-relayer/ethereum"
	"github.com/strangelove-ventures/noble-cctp-relayer/noble"
	"github.com/strangelove-ventures/noble-cctp-relayer/solana"
	"github.com/strangelove-ventures/noble-cctp-relayer/types"
)

var logger log.Logger
var EnvFile = os.ExpandEnv("$GOPATH/src/github.com/strangelove-ventures/noble-cctp-relayer/.env")
var LocalEnvFile = ".env"

func init() {
	// define logger
	logger = log.NewLogger(os.Stdout, log.LevelOption(zerolog.ErrorLevel))

	// Try loading a local .env first (project root), then GOPATH fallback, then relative from test dirs
	if err := godotenv.Load(LocalEnvFile); err != nil {
		// GOPATH-based fallback
		if err2 := godotenv.Load(EnvFile); err2 != nil {
			// Relative fallback when tests run from sub-packages like cmd/
			if err3 := godotenv.Load("../.env"); err3 != nil {
				// Log and continue; not all tests require .env
				wd, _ := os.Getwd()
				logger.Error(".env not found; continuing without env overrides", "local_err", err, "gopath_err", err2, "relative_err", err3, "cwd", wd)
			}
		}
	}
}

func ConfigSetup(t *testing.T) (a *cmd.AppState, registeredDomains map[types.Domain]types.Chain) {
	t.Helper()

	// Generate a throwaway Solana key for tests (no funds required)
	key, _ := extsolana.NewRandomPrivateKey()

	var testConfig = types.Config{
		Chains: map[string]types.ChainConfig{
			"noble": &noble.ChainConfig{
				ChainID: "grand-1",
				RPC:     os.Getenv("NOBLE_RPC"),
			},
			"ethereum": &ethereum.ChainConfig{
				ChainID:          11155111,
				Domain:           types.Domain(0),
				MinterPrivateKey: "1111111111111111111111111111111111111111111111111111111111111111",
				RPC:              os.Getenv("SEPOLIA_RPC"),
				WS:               os.Getenv("SEPOLIA_WS"),
			},
			"solana": &solana.Config{
				RPC:                  os.Getenv("SOLANA_RPC"),
				WS:                   os.Getenv("SOLANA_WS"),
				Domain:               5,
				MessageTransmitter:   "CCTPmbSD7gX1bxKPAmg77w8oFzNFpaQiQUWD43TKaecd",
				TokenMessengerMinter: "CCTPiPYPc6AsJuwueEnWgSgucamXDZwBd53dQ11YiKX3",
				FiatToken:            "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				RemoteTokens: map[types.Domain]string{
					0: "0x487039debedbf32d260137b0a6f66b90962bec777250910d253781de326a716d",
					4: "0x487039debedbf32d260137b0a6f66b90962bec777250910d253781de326a716d",
				},
				BroadcastRetries:       1,
				BroadcastRetryInterval: 1,
				MinMintAmount:          1,
				MinterPrivateKey:       key.String(),
			},
		},
		Circle: types.CircleSettings{
			AttestationBaseURL: "https://iris-api-sandbox.circle.com/attestations/",
			FetchRetries:       0,
			FetchRetryInterval: 3,
		},

		EnabledRoutes: map[types.Domain][]types.Domain{
			0: {4},
			4: {0},
		},
	}

	a = cmd.NewAppState()
	a.LogLevel = "debug"
	a.InitLogger()
	a.Config = &testConfig

	registeredDomains = make(map[types.Domain]types.Chain)
	for name, cfgg := range a.Config.Chains {
		c, err := cfgg.Chain(name)
		require.NoError(t, err, "Error creating chain")

		registeredDomains[c.Domain()] = c
	}

	return a, registeredDomains
}
