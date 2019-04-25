// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/moov-io/base"
	"github.com/moov-io/gl"
)

type testQLedgerAccountRepository struct {
	*qledgerAccountRepository

	deployment *qledgerDeployment
}

func (q *testQLedgerAccountRepository) close() {
	if q.deployment != nil {
		q.deployment.close()
	}
}

// qualifyQLedgerAccountTest will skip tests if Go's test -short flag is specified or if
// the needed env variables aren't set. See above for the env variables.
//
// Returned will be a qledgerAccountRepository
func qualifyQLedgerAccountTest(t *testing.T) *testQLedgerAccountRepository {
	t.Helper()

	if testing.Short() {
		t.Skip("-short flag enabled")
	}

	deployment := spawnQLedger(t)
	if deployment == nil {
		t.Fatal("nil QLedger docker deployment")
	}

	// repo, err := setupQLedgerAccountStorage("https://api.moov.io/v1/qledger", "moov") // Test against Production
	repo, err := setupQLedgerAccountStorage(fmt.Sprintf("http://localhost:%s", deployment.qledger.GetPort("7000/tcp")), "moov")
	if err != nil {
		t.Fatal(err)
	}
	return &testQLedgerAccountRepository{repo, deployment}
}

func TestQLedgerAccounts__ping(t *testing.T) {
	repo := qualifyQLedgerAccountTest(t)
	defer repo.close()

	if err := repo.Ping(); err != nil {
		t.Error(err)
	}
}

func TestQLedger__Accounts(t *testing.T) {
	repo := qualifyQLedgerAccountTest(t)

	customerId := base.ID()
	account := &gl.Account{
		ID:               base.ID(),
		CustomerID:       customerId,
		Name:             "example account",
		AccountNumber:    createAccountNumber(),
		RoutingNumber:    "121042882",
		Status:           "Active",
		Type:             "Checking",
		Balance:          100,
		BalancePending:   123,
		BalanceAvailable: 412,
		CreatedAt:        time.Now(),
	}
	if err := repo.CreateAccount(customerId, account); err != nil {
		t.Error(err)
	}

	// Now grab accounts for this customer
	accounts, err := repo.GetCustomerAccounts(customerId)
	if err != nil {
		t.Error(err)
	}
	if len(accounts) == 0 {
		t.Fatal("no accounts found")
	}
	if account.ID != accounts[0].ID {
		t.Errorf("expected account %q, but found %#v", account.ID, accounts[0].ID)
	}
	if account.Balance != 100 || account.BalancePending != 123 || account.BalanceAvailable != 412 {
		t.Errorf("Balance=%d BalancePending=%d BalanceAvailable=%d", account.Balance, account.BalancePending, account.BalanceAvailable)
	}
	if account.CreatedAt.IsZero() {
		t.Error("zero time for CreatedAt")
	}

	// Search for account
	acct, err := repo.SearchAccounts(account.AccountNumber, account.RoutingNumber, "Checking")
	if err != nil {
		t.Fatal(err)
	}
	if acct == nil {
		t.Fatal("SearchAccounts: nil account")
	}
	if acct.ID != account.ID {
		t.Errorf("acct.ID=%q account.ID=%q", acct.ID, account.ID)
	}

	repo.close()
}

func TestQLedger__read(t *testing.T) {
	if v := readBalance("100"); v != 100 {
		t.Errorf("got %v", v)
	}
	if v := readBalance("asas"); v != 0 {
		t.Errorf("got %v", v)
	}

	if v := readTime("2019-01-02T15:04:05Z").Format(time.RFC3339); v != "2019-01-02T15:04:05Z" {
		t.Errorf("got %q", v)
	}
}
