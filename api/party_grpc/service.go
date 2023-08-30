package party_grpc

import (
	"context"
	"fmt"

	"github.com/mrbanja/watchparty/tools/logging"

	"github.com/mrbanja/watchparty/party"

	"github.com/bufbuild/connect-go"
	gen "github.com/mrbanja/watchparty-proto/gen-go"
	"github.com/mrbanja/watchparty-proto/gen-go/protocolconnect"
	"go.uber.org/zap"
)

type service struct {
	logger *zap.Logger
	party  *party.Party
}

func New(p *party.Party, logger *zap.Logger) protocolconnect.PartyServiceHandler {
	return &service{
		logger: logger.Named("party_grpc"),
		party:  p,
	}
}

func (s service) GetMagnet(ctx context.Context, c *connect.Request[gen.GetMagnetRequest]) (*connect.Response[gen.GetMagnetResponse], error) {
	return connect.NewResponse(&gen.GetMagnetResponse{Magnet: "magnet:?..."}), nil
}

func (s service) JoinRoom(ctx context.Context, c *connect.BidiStream[gen.RoomRequest, gen.RoomResponse]) error {
	logger := logging.WithTrace(s.logger, ctx)
	msg, err := c.Receive()
	if err != nil {
		logger.Warn("Error while receiving initial join req", zap.Error(err))
		return err
	}
	conn := msg.GetConnect()
	if conn == nil {
		logger.Warn("Not initial message send during first join request")
		return fmt.Errorf("initial message should be `Connect`")
	}

	room := s.party.GetOrCreateRoom(conn.RoomName)
	participant := room.AddParticipant(ctx, c)
	defer room.RemoveParticipant(ctx, participant)

	for {
		if ctx.Err() != nil {
			logger.Warn("Context canceled")
			return nil
		}

		msg, err := c.Receive()
		if err != nil {
			logger.Warn("Error while receiving an update", zap.Error(err))
			return err
		}
		update := msg.GetUpdate()
		if update == nil {
			logger.Warn("Not update message send during join request")
			return fmt.Errorf("update should be sent after join")
		}

		room.BroadcastExcept(
			ctx,
			&gen.RoomResponse{Update: &gen.Update{
				State: update.State,
				Time:  update.Time,
			}},
			participant,
		)
	}
}
