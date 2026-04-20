package checkpoint

import (
	"encoding/json"
	"os"
	"sync"
)

// State holds the persisted checkpoint data.
type State struct {
	Offset int64  `json:"offset"`
	File   string `json:"file"`
}

// Checkpoint persists and restores stream read positions.
type Checkpoint struct {
	mu   sync.Mutex
	path string
	state State
}

// New creates a Checkpoint backed by the given file path.
func New(path string) (*Checkpoint, error) {
	c := &Checkpoint{path: path}
	if err := c.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return c, nil
}

// Save persists the current offset and source file to disk.
func (c *Checkpoint) Save(file string, offset int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = State{File: file, Offset: offset}
	data, err := json.Marshal(c.state)
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0644)
}

// Load returns the last saved State.
func (c *Checkpoint) Load() State {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

// Reset clears the persisted checkpoint file and in-memory state.
func (c *Checkpoint) Reset() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = State{}
	return os.Remove(c.path)
}

func (c *Checkpoint) load() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &c.state)
}
