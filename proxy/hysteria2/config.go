package hysteria2

import (
	"github.com/gzjjjfree/v5-result/common/protocol"
)

// MemoryAccount is an account type converted from Account.
type MemoryAccount struct{}

// AsAccount implements protocol.AsAccount.
func (a *Account) AsAccount() (protocol.Account, error) {
	return &MemoryAccount{}, nil
}

// Equals implements protocol.Account.Equals().
func (a *MemoryAccount) Equals(another protocol.Account) bool {
	return false
}
