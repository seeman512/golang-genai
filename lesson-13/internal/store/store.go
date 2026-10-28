// Package store contains the small in-memory domain store used by the lesson.
package store

import (
	"errors"
	"fmt"
	"sync"
	"time"

	pb "mystic-square-grpc/gen/mysticsquare"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrDuplicateEmail = errors.New("email already registered")
	ErrInvalidState   = errors.New("invalid registration state transition")
	ErrInvalidMove    = errors.New("invalid game move")
	ErrGameFinished   = errors.New("game is already finished")
)

type RegistrationEvent struct {
	RegistrationID uint64
	Status         pb.RegistrationStatus
	At             time.Time
	Message        string
	Sequence       int64
}

type Registration struct {
	ID        uint64
	Email     string
	FullName  string
	Status    pb.RegistrationStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Move struct{ Row, Column int }

type Game struct {
	ID      uint64
	Rows    int
	Columns int
	Board   []int
	Moves   []Move
	IsWin   bool
	IsOver  bool
}

type Store struct {
	mu            sync.RWMutex
	nextRegID     uint64
	nextGameID    uint64
	registrations map[uint64]Registration
	byEmail       map[string]uint64
	events        map[uint64][]RegistrationEvent
	subs          map[uint64]map[chan RegistrationEvent]struct{}
	games         map[uint64]Game
}

func New() *Store {
	return &Store{
		nextRegID: 1, nextGameID: 1,
		registrations: make(map[uint64]Registration), byEmail: make(map[string]uint64),
		events: make(map[uint64][]RegistrationEvent), subs: make(map[uint64]map[chan RegistrationEvent]struct{}),
		games: make(map[uint64]Game),
	}
}

func (s *Store) CreateRegistration(email, fullName string) (Registration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byEmail[email]; ok {
		return Registration{}, ErrDuplicateEmail
	}
	now := time.Now().UTC()
	reg := Registration{ID: s.nextRegID, Email: email, FullName: fullName, Status: pb.RegistrationStatus_PENDING, CreatedAt: now, UpdatedAt: now}
	s.nextRegID++
	s.registrations[reg.ID] = reg
	s.byEmail[email] = reg.ID
	s.appendEventLocked(reg.ID, pb.RegistrationStatus_UNREGISTERED, "registration request received", now)
	s.appendEventLocked(reg.ID, pb.RegistrationStatus_PENDING, "waiting for mail confirmation", now)
	return cloneRegistration(reg), nil
}

func (s *Store) GetRegistration(id uint64) (Registration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	reg, ok := s.registrations[id]
	if !ok {
		return Registration{}, ErrNotFound
	}
	return cloneRegistration(reg), nil
}

func (s *Store) AdvanceRegistration(id uint64, want pb.RegistrationStatus) (Registration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	reg, ok := s.registrations[id]
	if !ok {
		return Registration{}, ErrNotFound
	}
	valid := (reg.Status == pb.RegistrationStatus_PENDING && want == pb.RegistrationStatus_MAIL_CONFIRMED) ||
		(reg.Status == pb.RegistrationStatus_MAIL_CONFIRMED && want == pb.RegistrationStatus_REGISTERED)
	if !valid {
		return Registration{}, fmt.Errorf("%w: %s -> %s", ErrInvalidState, reg.Status, want)
	}
	reg.Status = want
	reg.UpdatedAt = time.Now().UTC()
	s.registrations[id] = reg
	message := "mail confirmed"
	if want == pb.RegistrationStatus_REGISTERED {
		message = "registration completed"
	}
	s.appendEventLocked(id, want, message, reg.UpdatedAt)
	return cloneRegistration(reg), nil
}

func (s *Store) Events(id uint64) ([]RegistrationEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.registrations[id]; !ok {
		return nil, ErrNotFound
	}
	return append([]RegistrationEvent(nil), s.events[id]...), nil
}

func (s *Store) Subscribe(id uint64) (<-chan RegistrationEvent, func()) {
	ch := make(chan RegistrationEvent, 16)
	s.mu.Lock()
	if s.subs[id] == nil {
		s.subs[id] = make(map[chan RegistrationEvent]struct{})
	}
	s.subs[id][ch] = struct{}{}
	s.mu.Unlock()
	var once sync.Once
	cancel := func() {
		once.Do(func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			if set := s.subs[id]; set != nil {
				delete(set, ch)
				close(ch)
				if len(set) == 0 {
					delete(s.subs, id)
				}
			}
		})
	}
	return ch, cancel
}

func (s *Store) SubscriberCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := 0
	for _, set := range s.subs {
		n += len(set)
	}
	return n
}

func (s *Store) CreateGame(rows, columns int) (Game, error) {
	if rows == 0 {
		rows = 2
	}
	if columns == 0 {
		columns = 2
	}
	if rows < 2 || columns < 2 || rows > 10 || columns > 10 {
		return Game{}, errors.New("rows and columns must be between 2 and 10")
	}
	board := make([]int, rows*columns)
	for i := range board {
		board[i] = i + 1
	}
	board[len(board)-1] = 0
	// Start one legal move away from solved, so the stream demo can win.
	empty := len(board) - 1
	board[empty], board[empty-1] = board[empty-1], board[empty]
	s.mu.Lock()
	defer s.mu.Unlock()
	game := Game{ID: s.nextGameID, Rows: rows, Columns: columns, Board: board}
	s.nextGameID++
	s.games[game.ID] = cloneGame(game)
	return cloneGame(game), nil
}

func (s *Store) GetGame(id uint64) (Game, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	game, ok := s.games[id]
	if !ok {
		return Game{}, ErrNotFound
	}
	return cloneGame(game), nil
}

func (s *Store) ApplyMove(id uint64, move Move) (Game, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	game, ok := s.games[id]
	if !ok {
		return Game{}, ErrNotFound
	}
	if game.IsWin || game.IsOver {
		return Game{}, ErrGameFinished
	}
	if move.Row < 0 || move.Row >= game.Rows || move.Column < 0 || move.Column >= game.Columns {
		return Game{}, ErrInvalidMove
	}
	empty := findEmpty(game.Board)
	index := move.Row*game.Columns + move.Column
	if abs(move.Row-empty/game.Columns)+abs(move.Column-empty%game.Columns) != 1 {
		return Game{}, ErrInvalidMove
	}
	game.Board[index], game.Board[empty] = game.Board[empty], game.Board[index]
	game.Moves = append(game.Moves, move)
	game.IsWin = solved(game.Board)
	if !game.IsWin && len(game.Moves) >= 100 {
		game.IsOver = true
	}
	s.games[id] = cloneGame(game)
	return cloneGame(game), nil
}

func (s *Store) appendEventLocked(id uint64, status pb.RegistrationStatus, message string, at time.Time) {
	event := RegistrationEvent{RegistrationID: id, Status: status, At: at, Message: message, Sequence: int64(len(s.events[id]) + 1)}
	s.events[id] = append(s.events[id], event)
	for ch := range s.subs[id] {
		select {
		case ch <- event:
		default:
		}
	}
}
func cloneRegistration(reg Registration) Registration { return reg }
func cloneGame(game Game) Game {
	game.Board = append([]int(nil), game.Board...)
	game.Moves = append([]Move(nil), game.Moves...)
	return game
}
func findEmpty(board []int) int {
	for i, v := range board {
		if v == 0 {
			return i
		}
	}
	return -1
}
func solved(board []int) bool {
	for i, v := range board {
		want := i + 1
		if i == len(board)-1 {
			want = 0
		}
		if v != want {
			return false
		}
	}
	return true
}
func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
