package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	pb "mystic-square-grpc/gen/mysticsquare"
)

func main() {
	addr := flag.String("addr", "localhost:50051", "gRPC server address")
	max := flag.Int("max", 0, "maximum status events to receive; zero means until Ctrl+C")
	flag.Parse()
	if len(flag.Args()) == 0 {
		usage()
		os.Exit(2)
	}
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fatal(err)
	}
	defer conn.Close()
	client := pb.NewMysticSquareServiceClient(conn)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	args := flag.Args()
	var runErr error
	switch args[0] {
	case "register":
		if len(args) != 3 {
			runErr = fmt.Errorf("usage: register <email> <full_name>")
		} else {
			runErr = register(ctx, client, args[1], args[2])
		}
	case "confirm":
		runErr = registrationAction(ctx, client, args, "confirm")
	case "complete":
		runErr = registrationAction(ctx, client, args, "complete")
	case "watch":
		if len(args) != 2 {
			runErr = fmt.Errorf("usage: watch <registration_id>")
		} else {
			runErr = watch(ctx, client, args[1], *max)
		}
	case "create-game":
		runErr = createGame(ctx, client, args)
	case "get-game":
		if len(args) != 2 {
			runErr = fmt.Errorf("usage: get-game <game_id>")
		} else {
			runErr = getGame(ctx, client, args[1])
		}
	case "moves":
		runErr = moves(ctx, client, args)
	default:
		usage()
		os.Exit(2)
	}
	if runErr != nil {
		fatal(runErr)
	}
}

func register(ctx context.Context, c pb.MysticSquareServiceClient, email, name string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	reg, err := c.Register(ctx, &pb.RegisterRequest{Email: email, FullName: name})
	if err != nil {
		return err
	}
	printMessage(reg)
	return nil
}
func registrationAction(ctx context.Context, c pb.MysticSquareServiceClient, args []string, kind string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: %s <registration_id>", kind)
	}
	id, err := parseUint(args[1])
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var reg *pb.Registration
	if kind == "confirm" {
		reg, err = c.ConfirmRegistration(ctx, &pb.RegistrationRequest{RegistrationId: id})
	} else {
		reg, err = c.CompleteRegistration(ctx, &pb.RegistrationRequest{RegistrationId: id})
	}
	if err != nil {
		return err
	}
	printMessage(reg)
	return nil
}
func watch(ctx context.Context, c pb.MysticSquareServiceClient, idText string, max int) error {
	id, err := parseUint(idText)
	if err != nil {
		return err
	}
	stream, err := c.WatchRegistrationStatus(ctx, &pb.WatchRegistrationRequest{RegistrationId: id, MaxEvents: int32(max)})
	if err != nil {
		return err
	}
	for {
		event, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			fmt.Println("status stream completed")
			return nil
		}
		if recvErr != nil {
			return recvErr
		}
		printMessage(event)
	}
}
func createGame(ctx context.Context, c pb.MysticSquareServiceClient, args []string) error {
	if len(args) > 3 {
		return fmt.Errorf("usage: create-game [rows columns]")
	}
	rows, columns := int32(2), int32(2)
	var err error
	if len(args) == 3 {
		rows, err = parseInt32(args[1])
		if err != nil {
			return err
		}
		columns, err = parseInt32(args[2])
		if err != nil {
			return err
		}
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	game, err := c.CreateGame(ctx, &pb.CreateGameRequest{Rows: rows, Columns: columns})
	if err != nil {
		return err
	}
	printMessage(game)
	return nil
}
func getGame(ctx context.Context, c pb.MysticSquareServiceClient, text string) error {
	id, err := parseUint(text)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	game, err := c.GetGame(ctx, &pb.GetGameRequest{GameId: id})
	if err != nil {
		return err
	}
	printMessage(game)
	return nil
}
func moves(ctx context.Context, c pb.MysticSquareServiceClient, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: moves <game_id> <row,column>...")
	}
	gameID, err := parseUint(args[1])
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	stream, err := c.SubmitGameMoves(ctx)
	if err != nil {
		return err
	}
	for i, text := range args[2:] {
		parts := strings.Split(text, ",")
		if len(parts) != 2 {
			return fmt.Errorf("move %q must be row,column", text)
		}
		row, e := parseInt32(parts[0])
		if e != nil {
			return e
		}
		column, e := parseInt32(parts[1])
		if e != nil {
			return e
		}
		if err := stream.Send(&pb.GameMove{GameId: gameID, Row: row, Column: column, Sequence: int64(i + 1)}); err != nil {
			return err
		}
	}
	summary, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	printMessage(summary)
	return nil
}
func parseUint(s string) (uint64, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil || v == 0 {
		return 0, fmt.Errorf("invalid positive id %q", s)
	}
	return v, nil
}
func parseInt32(s string) (int32, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", s)
	}
	return int32(v), nil
}
func printMessage(message proto.Message) {
	out, err := protojson.MarshalOptions{Multiline: true, UseProtoNames: true, EmitUnpopulated: true}.Marshal(message)
	if err != nil {
		fatal(err)
	}
	fmt.Println(string(out))
}
func fatal(err error) { fmt.Fprintln(os.Stderr, "error:", err); os.Exit(1) }
func usage() {
	fmt.Fprint(os.Stderr, `mystic-client [flags] <command>

Commands:
  register <email> <full_name>
  confirm <registration_id>
  complete <registration_id>
  watch <registration_id>
  create-game [rows columns]
  get-game <game_id>
  moves <game_id> <row,column>...
`)
	flag.PrintDefaults()
}
