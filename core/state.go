package core

import (
	"errors"
	"sync"
)

type State struct {
	mu       sync.RWMutex
	Names    map[string]string
	Nonces   map[string]uint64
	Accounts map[string]bool // Tracks addresses we've seen before
}

func NewState() *State {
	return &State{
		Names:    make(map[string]string),
		Nonces:   make(map[string]uint64),
		Accounts: make(map[string]bool),
	}
}

func (s *State) ApplyTransaction(tx Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// First-time address setup
	if !s.Accounts[tx.Address] {
		s.Accounts[tx.Address] = true
		s.Nonces[tx.Address] = 0
	}

	err := s.validateTransactionWithoutLock(tx)
	if err != nil {
		return err
	}

	// Update state
	s.Names[tx.Name] = tx.Address
	s.Nonces[tx.Address]++
	return nil
}

func (s *State) GetAddressByName(name string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	addr, ok := s.Names[name]
	return addr, ok
}

func (s *State) GetNonce(address string) uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Nonces[address]
}

func (s *State) ValidateTransaction(tx Transaction) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.validateTransactionWithoutLock(tx)
}

func (s *State) validateTransactionWithoutLock(tx Transaction) error {
	// Verify signature and address match
	if err := tx.VerifySignature(); err != nil {
		return err
	}

	// Check nonce
	if s.Accounts[tx.Address] && tx.Nonce != s.Nonces[tx.Address] {
		return errors.New("invalid nonce")
	}

	// Check if name is already registered
	if _, exists := s.Names[tx.Name]; exists {
		return errors.New("name already registered")
	}

	return nil
}
