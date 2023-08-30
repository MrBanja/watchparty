package app

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/mrbanja/watchparty/tools/logging"

	"github.com/mrbanja/watchparty/party"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/mrbanja/watchparty-proto/gen-go/protocolconnect"
	"github.com/mrbanja/watchparty/api/party_grpc"
	"github.com/mrbanja/watchparty/api/ssr"
	"github.com/mrbanja/watchparty/tools/http_server"
)

type Options struct {
	PublicAddr string `env:"PUBLIC_ADDR,required"`
	LocalAddr  string `env:"LOCAL_ADDR" envDefault:":8000"`
}

func Run(ctx context.Context, opt Options) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c := zap.NewDevelopmentConfig()
	c.Development = false
	logger, err := c.Build()
	if err != nil {
		return err
	}
	logger = logger.Named("app")

	prt := party.New(logger)
	ssrSRV := ssr.New(prt, opt.PublicAddr, logger)
	ssrSRV.MustBuildTemplate()
	partyGRPC := party_grpc.New(prt, logger)

	api := mux.NewRouter()
	api.HandleFunc("/{room_id}/peer_status/{participant_id}", ssrSRV.GetStatusBadge).Methods(http.MethodGet)

	prefix, partyGRPCHandler := protocolconnect.NewPartyServiceHandler(partyGRPC)
	api.PathPrefix(prefix).Handler(partyGRPCHandler)

	handler := logging.Middleware(api, logger)

	server := &http.Server{
		Addr: opt.PublicAddr,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		Handler: h2c.NewHandler(handler, &http2.Server{}),
	}

	logger.Info("[*] Http server started", zap.String("Pub Addr", opt.PublicAddr), zap.String("Local Addr", opt.LocalAddr))
	if err = http_server.Serve(ctx, server, 10*time.Second, func(server *http.Server) error {
		return server.ListenAndServe()
	}); err != nil {
		return err
	}

	return nil
}
