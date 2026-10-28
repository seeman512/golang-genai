package server

import (
	"context"
	"errors"
	"io"
	"net/mail"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "mystic-square-grpc/gen/mysticsquare"
	"mystic-square-grpc/internal/store"
)

// Server implements the Mystic Square protobuf service.
type Server struct {
	pb.UnimplementedMysticSquareServiceServer
	store *store.Store
}

// New creates a gRPC server backed by st.
func New(st *store.Store) *Server { return &Server{store: st} }

func (s *Server) Register(_ context.Context, req *pb.RegisterRequest) (*pb.Registration, error) {
	email := strings.ToLower(strings.TrimSpace(req.GetEmail()))
	fullName := strings.TrimSpace(req.GetFullName())
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email || fullName == "" {
		return nil, status.Error(codes.InvalidArgument, "valid email and full_name are required")
	}
	reg, err := s.store.CreateRegistration(email, fullName)
	if err != nil {
		return nil, toStatus(err)
	}
	return registrationMessage(reg), nil
}

func (s *Server) ConfirmRegistration(_ context.Context, req *pb.RegistrationRequest) (*pb.Registration, error) {
	if req.GetRegistrationId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "registration_id is required")
	}
	reg, err := s.store.AdvanceRegistration(req.GetRegistrationId(), pb.RegistrationStatus_MAIL_CONFIRMED)
	if err != nil {
		return nil, toStatus(err)
	}
	return registrationMessage(reg), nil
}

func (s *Server) CompleteRegistration(_ context.Context, req *pb.RegistrationRequest) (*pb.Registration, error) {
	if req.GetRegistrationId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "registration_id is required")
	}
	reg, err := s.store.AdvanceRegistration(req.GetRegistrationId(), pb.RegistrationStatus_REGISTERED)
	if err != nil {
		return nil, toStatus(err)
	}
	return registrationMessage(reg), nil
}

func (s *Server) WatchRegistrationStatus(req *pb.WatchRegistrationRequest, stream pb.MysticSquareService_WatchRegistrationStatusServer) error {
	id := req.GetRegistrationId()
	if id == 0 {
		return status.Error(codes.InvalidArgument, "registration_id is required")
	}
	if _, err := s.store.GetRegistration(id); err != nil {
		return toStatus(err)
	}
	changes, cancel := s.store.Subscribe(id)
	defer cancel()
	history, err := s.store.Events(id)
	if err != nil {
		return toStatus(err)
	}
	max := req.GetMaxEvents()
	sent := int32(0)
	lastSequence := int64(0)
	send := func(event store.RegistrationEvent) error {
		if max > 0 && sent >= max {
			return nil
		}
		if err := stream.Send(&pb.RegistrationStatusEvent{RegistrationId: event.RegistrationID, Status: event.Status, OccurredAt: timestamppb.New(event.At), Message: event.Message, Sequence: event.Sequence}); err != nil {
			return err
		}
		sent++
		lastSequence = event.Sequence
		return nil
	}
	for _, event := range history {
		if err := send(event); err != nil {
			return err
		}
		if max > 0 && sent >= max {
			return nil
		}
	}
	for {
		select {
		case <-stream.Context().Done():
			return status.FromContextError(stream.Context().Err()).Err()
		case event, ok := <-changes:
			if !ok {
				return nil
			}
			if event.Sequence <= lastSequence {
				continue
			}
			if err := send(event); err != nil {
				return err
			}
			if max > 0 && sent >= max {
				return nil
			}
		}
	}
}

func (s *Server) CreateGame(_ context.Context, req *pb.CreateGameRequest) (*pb.GameSnapshot, error) {
	game, err := s.store.CreateGame(int(req.GetRows()), int(req.GetColumns()))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return gameMessage(game), nil
}

func (s *Server) GetGame(_ context.Context, req *pb.GetGameRequest) (*pb.GameSnapshot, error) {
	if req.GetGameId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "game_id is required")
	}
	game, err := s.store.GetGame(req.GetGameId())
	if err != nil {
		return nil, toStatus(err)
	}
	return gameMessage(game), nil
}

func (s *Server) SubmitGameMoves(stream pb.MysticSquareService_SubmitGameMovesServer) error {
	summary := &pb.GameMoveSummary{}
	var gameID uint64
	for {
		move, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			if gameID == 0 {
				return status.Error(codes.InvalidArgument, "at least one game move is required")
			}
			game, getErr := s.store.GetGame(gameID)
			if getErr != nil {
				return toStatus(getErr)
			}
			summary.GameId = gameID
			summary.IsWin = game.IsWin
			summary.IsOver = game.IsOver
			summary.Board = int32Board(game.Board)
			return stream.SendAndClose(summary)
		}
		if err != nil {
			return err
		}
		summary.Received++
		if move.GetGameId() == 0 {
			return status.Error(codes.InvalidArgument, "game_id is required")
		}
		if gameID == 0 {
			gameID = move.GetGameId()
		} else if gameID != move.GetGameId() {
			return status.Error(codes.InvalidArgument, "all streamed moves must use the same game_id")
		}
		game, applyErr := s.store.ApplyMove(gameID, store.Move{Row: int(move.GetRow()), Column: int(move.GetColumn())})
		if applyErr != nil {
			summary.Rejected++
			if len(summary.Rejections) < 32 {
				summary.Rejections = append(summary.Rejections, &pb.MoveRejection{Sequence: move.GetSequence(), Reason: applyErr.Error()})
			}
			continue
		}
		summary.Applied++
		summary.IsWin = game.IsWin
		summary.IsOver = game.IsOver
		summary.Board = int32Board(game.Board)
	}
}

func registrationMessage(reg store.Registration) *pb.Registration {
	return &pb.Registration{Id: reg.ID, Email: reg.Email, FullName: reg.FullName, Status: reg.Status, CreatedAt: timestamppb.New(reg.CreatedAt), UpdatedAt: timestamppb.New(reg.UpdatedAt)}
}
func gameMessage(game store.Game) *pb.GameSnapshot {
	return &pb.GameSnapshot{Id: game.ID, Board: int32Board(game.Board), Moves: int32(len(game.Moves)), IsWin: game.IsWin, IsOver: game.IsOver}
}
func int32Board(board []int) []int32 {
	out := make([]int32, len(board))
	for i, v := range board {
		out[i] = int32(v)
	}
	return out
}
func toStatus(err error) error {
	switch {
	case errors.Is(err, store.ErrNotFound):
		return status.Error(codes.NotFound, "resource not found")
	case errors.Is(err, store.ErrDuplicateEmail):
		return status.Error(codes.AlreadyExists, "email already registered")
	case errors.Is(err, store.ErrInvalidState):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, store.ErrInvalidMove):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, store.ErrGameFinished):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
