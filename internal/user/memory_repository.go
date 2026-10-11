package user

import "sync"

type MemoryRepository struct {
	mu           sync.RWMutex
	nextID       int64
	byID         map[int64]User
	idByNickname map[string]int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextID:       1,
		byID:         make(map[int64]User),
		idByNickname: make(map[string]int64),
	}
}

func (r *MemoryRepository) Create(u User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := u.Nickname
	if _, exists := r.idByNickname[key]; exists {
		return User{}, ErrNicknameTaken
	}

	u.ID = r.nextID
	r.nextID++

	r.byID[u.ID] = u
	r.idByNickname[key] = u.ID

	return u, nil
}

func (r *MemoryRepository) GetByID(id int64) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, exists := r.byID[id]
	if !exists {
		return User{}, ErrNotFound
	}

	return u, nil
}

func (r *MemoryRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, exists := r.byID[id]
	if !exists {
		return ErrNotFound
	}

	delete(r.byID, id)
	delete(r.idByNickname, u.Nickname)

	return nil
}
