package client_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/stretchr/testify/assert"
)

func newSendErrorWrapped(s string) *client.SendError {
	return client.NewSendError(errors.Wrap(errors.New(s), "wrapped with some old bollocks"))
}

func Test_Eth_Errors(t *testing.T) {
	t.Parallel()
	var err *client.SendError
	randomError := client.NewSendErrorS("some old bollocks")

	t.Run("IsNonceTooLowError", func(t *testing.T) {
		assert.False(t, randomError.IsNonceTooLowError())

		tests := []struct {
			message string
			expect  bool
		}{
			// Geth
			{"nonce too low", true},
			// Parity
			{"Transaction nonce is too low. Try incrementing the nonce.", true},
			// Arbitrum
			{"transaction rejected: nonce too low", true},
			{"invalid transaction nonce", true},
			// Optimism
			{"invalid transaction: nonce too low", true},
			// Avalanche
			{"call failed: nonce too low: address 0x0499BEA33347cb62D79A9C0b1EDA01d8d329894c current nonce (5833) > tx nonce (5511)", true},
		}

		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsNonceTooLowError(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsNonceTooLowError(), test.expect)
		}
	})

	t.Run("IsReplacementUnderpriced", func(t *testing.T) {

		tests := []struct {
			message string
			expect  bool
		}{
			// Geth
			{"replacement transaction underpriced", true},
			// Parity
			{"Transaction gas price 100wei is too low. There is another transaction with same nonce in the queue with gas price 150wei. Try increasing the gas price or incrementing the nonce.", true},
			{"There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.", false},
			// Arbitrum
			{"gas price too low", false},
		}

		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsReplacementUnderpriced(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsReplacementUnderpriced(), test.expect)
		}
	})

	t.Run("IsTransactionAlreadyInMempool", func(t *testing.T) {
		assert.False(t, randomError.IsTransactionAlreadyInMempool())

		tests := []struct {
			message string
			expect  bool
		}{
			// Geth
			// I have seen this in log output
			{"known transaction: 0x7f657507aee0511e36d2d1972a6b22e917cc89f92b6c12c4dbd57eaabb236960", true},
			// This comes from the geth source - https://github.com/ethereum/go-ethereum/blob/eb9d7d15ecf08cd5104e01a8af64489f01f700b0/core/tx_pool.go#L57
			{"already known", true},
			// This one is present in the light client (?!)
			{"Known transaction (7f65)", true},
			// Parity
			{"Transaction with the same hash was already imported.", true},
		}
		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsTransactionAlreadyInMempool(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsTransactionAlreadyInMempool(), test.expect)
		}
	})

	t.Run("IsTerminallyUnderpriced", func(t *testing.T) {
		assert.False(t, randomError.IsTerminallyUnderpriced())

		tests := []struct {
			message string
			expect  bool
		}{
			// Geth
			{"transaction underpriced", true},
			{"replacement transaction underpriced", false},
			// Parity
			{"There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.", false},
			{"Transaction gas price is too low. It does not satisfy your node's minimal gas price (minimal: 100 got: 50). Try increasing the gas price.", true},
			// Arbitrum
			{"gas price too low", true},
		}

		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsTerminallyUnderpriced(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsTerminallyUnderpriced(), test.expect)
		}
	})

	t.Run("IsTemporarilyUnderpriced", func(t *testing.T) {
		tests := []struct {
			message string
			expect  bool
		}{
			// Parity
			{"There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.", true},
			{"There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.", true},
			{"Transaction gas price is too low. It does not satisfy your node's minimal gas price (minimal: 100 got: 50). Try increasing the gas price.", false},
		}
		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsTemporarilyUnderpriced(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsTemporarilyUnderpriced(), test.expect)
		}
	})

	t.Run("IsInsufficientEth", func(t *testing.T) {
		tests := []struct {
			message string
			expect  bool
		}{
			// Geth
			{"insufficient funds for transfer", true},
			{"insufficient funds for gas * price + value", true},
			{"insufficient balance for transfer", true},
			// Parity
			{"Insufficient balance for transaction. Balance=100.25, Cost=200.50", true},
			{"Insufficient funds. The account you tried to send transaction from does not have enough funds. Required 200.50 and got: 100.25.", true},
			// Arbitrum
			{"transaction rejected: insufficient funds for gas * price + value", true},
			{"not enough funds for gas", true},
			// Optimism
			{"invalid transaction: insufficient funds for gas * price + value", true},
		}
		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsInsufficientEth(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsInsufficientEth(), test.expect)
		}
	})

	t.Run("IsTooExpensive", func(t *testing.T) {
		tests := []struct {
			message string
			expect  bool
		}{
			// Geth
			{"tx fee (1.10 ether) exceeds the configured cap (1.00 ether)", true},
		}
		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsTooExpensive(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsTooExpensive(), test.expect)
		}

		assert.False(t, randomError.IsTooExpensive())
		// Nil
		err = client.NewSendError(nil)
		assert.False(t, err.IsTooExpensive())
	})

	t.Run("IsNonceMax", func(t *testing.T) {
		tests := []struct {
			message string
			expect  bool
		}{
			// Geth
			{"nonce has max value", true},
		}
		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsNonceMax(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsNonceMax(), test.expect)
		}
		assert.False(t, randomError.IsNonceMax())
		// Nil
		err = client.NewSendError(nil)
		assert.False(t, err.IsNonceMax())
	})

	t.Run("IsGasLimitReached", func(t *testing.T) {
		tests := []struct {
			message string
			expect  bool
		}{
			// Geth
			{"gas limit reached", true},
		}
		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsGasLimitReached(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsGasLimitReached(), test.expect)
		}
		assert.False(t, randomError.IsGasLimitReached())
		// Nil
		err = client.NewSendError(nil)
		assert.False(t, err.IsGasLimitReached())
	})

	t.Run("IsTipAboveFeeCap", func(t *testing.T) {
		tests := []struct {
			message string
			expect  bool
		}{
			// Geth
			{"max priority fee per gas higher than max fee per gas", true},
		}
		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsTipAboveFeeCap(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsTipAboveFeeCap(), test.expect)
		}
		assert.False(t, randomError.IsTipAboveFeeCap())
		// Nil
		err = client.NewSendError(nil)
		assert.False(t, err.IsTipAboveFeeCap())
	})

	t.Run("IsTipVeryHigh", func(t *testing.T) {
		tests := []struct {
			message string
			expect  bool
		}{
			// Geth
			{"max priority fee per gas higher than 2^256-1", true},
		}
		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsTipVeryHigh(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsTipVeryHigh(), test.expect)
		}
		assert.False(t, randomError.IsTipVeryHigh())
		// Nil
		err = client.NewSendError(nil)
		assert.False(t, err.IsTipVeryHigh())
	})

	t.Run("IsFeeCapVeryHigh", func(t *testing.T) {
		tests := []struct {
			message string
			expect  bool
		}{
			// Geth
			{"max fee per gas higher than 2^256-1", true},
		}
		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsFeeCapVeryHigh(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsFeeCapVeryHigh(), test.expect)
		}
		assert.False(t, randomError.IsFeeCapVeryHigh())
		// Nil
		err = client.NewSendError(nil)
		assert.False(t, err.IsFeeCapVeryHigh())
	})

	t.Run("IsFeeCapTooLow", func(t *testing.T) {
		tests := []struct {
			message string
			expect  bool
		}{
			// Geth
			{"max fee per gas less than block base fee", true},
		}
		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsFeeCapTooLow(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsFeeCapTooLow(), test.expect)
		}
		assert.False(t, randomError.IsFeeCapTooLow())
		// Nil
		err = client.NewSendError(nil)
		assert.False(t, err.IsFeeCapTooLow())
	})

	t.Run("IsSenderNotEOA", func(t *testing.T) {
		tests := []struct {
			message string
			expect  bool
		}{
			// Geth
			{"sender not an eoa", true},
		}
		for _, test := range tests {
			err = client.NewSendErrorS(test.message)
			assert.Equal(t, err.IsSenderNotEOA(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsSenderNotEOA(), test.expect)
		}
		assert.False(t, randomError.IsSenderNotEOA())
		// Nil
		err = client.NewSendError(nil)
		assert.False(t, err.IsSenderNotEOA())
	})

	t.Run("Optimism Fees errors", func(t *testing.T) {
		err := client.NewSendErrorS("primary websocket (wss://ws-mainnet.optimism.io) call failed: fee too high: 5835750750000000, use less than 467550750000000 * 0.700000")
		assert.True(t, err.IsFeeTooHigh())
		assert.False(t, err.IsFeeTooLow())
		err = newSendErrorWrapped("primary websocket (wss://ws-mainnet.optimism.io) call failed: fee too high: 5835750750000000, use less than 467550750000000 * 0.700000")
		assert.True(t, err.IsFeeTooHigh())
		assert.False(t, err.IsFeeTooLow())

		err = client.NewSendErrorS("fee too low: 30365610000000, use at least tx.gasLimit = 5874374 and tx.gasPrice = 15000000")
		assert.False(t, err.IsFeeTooHigh())
		assert.True(t, err.IsFeeTooLow())
		err = newSendErrorWrapped("fee too low: 30365610000000, use at least tx.gasLimit = 5874374 and tx.gasPrice = 15000000")
		assert.False(t, err.IsFeeTooHigh())
		assert.True(t, err.IsFeeTooLow())

		assert.False(t, randomError.IsFeeTooHigh())
		assert.False(t, randomError.IsFeeTooLow())
		// Nil
		err = client.NewSendError(nil)
		assert.False(t, err.IsFeeTooHigh())
		assert.False(t, err.IsFeeTooLow())
	})

	t.Run("moonriver errors", func(t *testing.T) {
		err := client.NewSendErrorS("primary http (http://***REDACTED***:9933) call failed: submit transaction to pool failed: Pool(Stale)")
		assert.True(t, err.IsNonceTooLowError())
		assert.False(t, err.IsTransactionAlreadyInMempool())
		assert.False(t, err.Fatal())
		err = client.NewSendErrorS("primary http (http://***REDACTED***:9933) call failed: submit transaction to pool failed: Pool(AlreadyImported)")
		assert.True(t, err.IsTransactionAlreadyInMempool())
		assert.False(t, err.IsNonceTooLowError())
		assert.False(t, err.Fatal())
	})
}

func Test_Eth_Errors_Fatal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		errStr      string
		expectFatal bool
	}{
		{"some old bollocks", false},

		// Geth
		{"insufficient funds for transfer", false},

		{"exceeds block gas limit", true},
		{"invalid sender", true},
		{"negative value", true},
		{"oversized data", true},
		{"gas uint64 overflow", true},
		{"intrinsic gas too low", true},
		{"nonce too high", true},

		// Parity
		{"Insufficient funds. The account you tried to send transaction from does not have enough funds. Required 100 and got: 50.", false},

		{"Supplied gas is beyond limit.", true},
		{"Sender is banned in local queue.", true},
		{"Recipient is banned in local queue.", true},
		{"Code is banned in local queue.", true},
		{"Transaction is not permitted.", true},
		{"Transaction is too big, see chain specification for the limit.", true},
		{"Transaction gas is too low. There is not enough gas to cover minimal cost of the transaction (minimal: 100 got: 50) Try increasing supplied gas.", true},
		{"Transaction cost exceeds current gas limit. Limit: 50, got: 100. Try decreasing supplied gas.", true},
		{"Invalid signature: some old bollocks", true},
		{"Invalid RLP data: some old bollocks", true},

		// Arbitrum
		{"invalid message format", true},
		{"forbidden sender address", true},
		{"tx dropped due to L2 congestion", false},
		{"execution reverted: error code", true},
	}

	for _, test := range tests {
		t.Run(test.errStr, func(t *testing.T) {
			err := client.NewSendError(errors.New(test.errStr))
			assert.Equal(t, test.expectFatal, err.Fatal())
		})
	}
}