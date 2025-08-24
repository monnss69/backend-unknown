package components

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

// Store wraps database operations for components.
type Store struct {
	db *sql.DB
}

// NewStore creates a component store.
func NewStore(db *sql.DB) *Store { return &Store{db: db} }

// Create inserts a new component.
func (s *Store) Create(ctx context.Context, c *Component) error {
	props, err := json.Marshal(c.PropsSchema)
	if err != nil {
		return fmt.Errorf("marshal props: %w", err)
	}
	const q = `INSERT INTO components (id, name, code, props_schema) VALUES ($1, $2, $3, $4)`
	_, err = s.db.ExecContext(ctx, q, c.ID, c.Name, c.Code, props)
	return err
}

// Get fetches a component by ID.
func (s *Store) Get(ctx context.Context, id string) (*Component, error) {
	const q = `SELECT id, name, code, props_schema FROM components WHERE id = $1`
	var (
		c   Component
		raw []byte
	)
	err := s.db.QueryRowContext(ctx, q, id).Scan(&c.ID, &c.Name, &c.Code, &raw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("query: %w", err)
	}
	if err := json.Unmarshal(raw, &c.PropsSchema); err != nil {
		return nil, fmt.Errorf("unmarshal props: %w", err)
	}
	return &c, nil
}

// List returns all components.
func (s *Store) List(ctx context.Context) ([]Component, error) {
	const q = `SELECT id, name, code, props_schema FROM components`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()
	var comps []Component
	for rows.Next() {
		var (
			c   Component
			raw []byte
		)
		if err := rows.Scan(&c.ID, &c.Name, &c.Code, &raw); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		if err := json.Unmarshal(raw, &c.PropsSchema); err != nil {
			return nil, fmt.Errorf("unmarshal props: %w", err)
		}
		comps = append(comps, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return comps, nil
}
