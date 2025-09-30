package controller

import (
	"context"
	"shantaram/app/api"
	"shantaram/app/service/auth"
	"shantaram/app/service/limits"
	"shantaram/app/service/menu"
	"shantaram/app/service/order"
	"shantaram/app/service/params"
	"shantaram/app/service/pubsub"
	"shantaram/pkg/config"
	"shantaram/pkg/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do"
)

var _ api.StrictServerInterface = (*Server)(nil)

type Server struct {
	appCtx        context.Context
	cfg           *config.Config
	dbConn        *pgxpool.Pool
	queries       *database.Queries
	authService   *auth.Service
	limitsService *limits.Service
	pubsubService *pubsub.Service
	menuService   *menu.Service
	orderService  *order.Service
	paramsService *params.Service
}

func NewStrictServer(di *do.Injector) *Server {
	return &Server{
		appCtx:        do.MustInvoke[context.Context](di),
		cfg:           do.MustInvoke[*config.Config](di),
		dbConn:        do.MustInvoke[*pgxpool.Pool](di),
		queries:       do.MustInvoke[*database.Queries](di),
		authService:   do.MustInvoke[*auth.Service](di),
		limitsService: do.MustInvoke[*limits.Service](di),
		pubsubService: do.MustInvoke[*pubsub.Service](di),
		menuService:   do.MustInvoke[*menu.Service](di),
		orderService:  do.MustInvoke[*order.Service](di),
		paramsService: do.MustInvoke[*params.Service](di),
	}
}
