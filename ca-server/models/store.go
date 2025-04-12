package models

import (
	"errors"
	"sync"
)

// Store defines the data access interface
type Store interface {
	GetUser(id string) (*User, error)
	ListUsers() ([]*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id string) error
}

// MemoryStore provides an in-memory implementation of Store
type MemoryStore struct {
	users  map[string]*User
	mutex  sync.RWMutex
	nextID int
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:  make(map[string]*User),
		nextID: 1,
	}
}

// GetUser retrieves a user by ID
func (s *MemoryStore) GetUser(id string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// ListUsers returns all users
func (s *MemoryStore) ListUsers() ([]*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	return users, nil
}

// CreateUser adds a new user
func (s *MemoryStore) CreateUser(user *User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Generate ID if not provided
	if user.ID == "" {
		user.ID = generateID(s.nextID)
		s.nextID++
	}

	s.users[user.ID] = user
	return nil
}

// UpdateUser updates an existing user
func (s *MemoryStore) UpdateUser(user *User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	s.users[user.ID] = user
	return nil
}

// DeleteUser removes a user
func (s *MemoryStore) DeleteUser(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(s.users, id)
	return nil
}

// Helper to generate a string ID
func generateID(id int) string {
	return string(rune('0' + id))
}
