package evilspammer

import (
	"testing"
	"time"

	"github.com/iotaledger/goshimmer/client/evilwallet"
	"github.com/stretchr/testify/require"
)

func TestSpamTransactions(t *testing.T) {
	evilWallet := evilwallet.NewEvilWallet()

	err := evilWallet.RequestFreshBigFaucetWallet()
	require.NoError(t, err)

	outWallet := evilWallet.NewWallet(evilwallet.Reuse)

	scenario := evilwallet.NewEvilScenario(evilwallet.SingleTransactionBatch(), false, outWallet)
	options := []Options{
		WithSpamRate(5, time.Second),
		WithBatchesSent(20),
		WithSpammingFunc(ValueSpammingFunc),
		WithSpamWallet(evilWallet),
		WithEvilScenario(scenario),
	}
	spammer := NewSpammer(options...)
	spammer.Spam()
}

func TestSpamDoubleSpend(t *testing.T) {
	evilWallet := evilwallet.NewEvilWallet()

	err := evilWallet.RequestFreshBigFaucetWallet()
	require.NoError(t, err)

	outWallet := evilWallet.NewWallet(evilwallet.Reuse)

	scenarioTx := evilwallet.NewEvilScenario(evilwallet.SingleTransactionBatch(), false, outWallet)
	scenarioDs := evilwallet.NewEvilScenario(evilwallet.DoubleSpendBatch(5), false, outWallet)
	customScenario := evilwallet.NewEvilScenario(evilwallet.Scenario1(), false, outWallet)
	options := []Options{
		WithSpamRate(5, time.Second),
		WithSpamDuration(time.Second * 10),
		WithSpammingFunc(CustomConflictSpammingFunc),
		WithSpamWallet(evilWallet),
	}
	txOptions := append(options, WithEvilScenario(scenarioTx))
	dsOptions := append(options, WithEvilScenario(scenarioDs))
	customOptions := append(options, WithEvilScenario(customScenario))

	txSpammer := NewSpammer(txOptions...)
	dsSpammer := NewSpammer(dsOptions...)
	customSpammer := NewSpammer(customOptions...)

	txSpammer.Spam()
	dsSpammer.Spam()
	customSpammer.Spam()

}