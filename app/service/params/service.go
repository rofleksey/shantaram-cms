package params

import (
	"context"
	"fmt"
	"log/slog"
	"shantaram/pkg/config"
	"shantaram/pkg/database"
	"shantaram/pkg/telemetry"
	"time"

	"github.com/rofleksey/meg"
	"github.com/samber/do"
)

var serviceName = "params"

type Service struct {
	cfg     *config.Config
	queries *database.Queries
	tracing *telemetry.Tracing
}

func New(di *do.Injector) (*Service, error) {
	return &Service{
		cfg:     do.MustInvoke[*config.Config](di),
		queries: do.MustInvoke[*database.Queries](di),
		tracing: do.MustInvoke[*telemetry.Tracing](di),
	}, nil
}

func (s *Service) SetHeaderText(ctx context.Context, text *string, deadline *time.Time) error {
	ctx, span := s.tracing.StartServiceSpan(ctx, serviceName, "set_header_text")
	defer span.End()

	if err := s.queries.SetParamsHeader(ctx, database.SetParamsHeaderParams{
		HeaderText:     text,
		HeaderDeadline: deadline,
	}); err != nil {
		return s.tracing.Error(span, fmt.Errorf("SetParamsHeader: %w", err))
	}

	s.tracing.Success(span)

	return nil
}

func (s *Service) RunHeaderDeadline(ctx context.Context) {
	meg.RunTicker(ctx, time.Minute, func() {
		if err := s.checkHeaderDeadline(ctx); err != nil {
			slog.Error("checkHeaderDeadline error",
				slog.Any("error", err),
			)
		}
	})
}

func (s *Service) checkHeaderDeadline(ctx context.Context) error {
	params, err := s.queries.GetParams(ctx)
	if err != nil {
		return fmt.Errorf("GetParams: %w", err)
	}

	if params.HeaderDeadline != nil && time.Now().After(*params.HeaderDeadline) {
		if err := s.queries.SetParamsHeader(ctx, database.SetParamsHeaderParams{
			HeaderText:     nil,
			HeaderDeadline: nil,
		}); err != nil {
			return fmt.Errorf("SetParamsHeader: %w", err)
		}

		slog.Info("Header was reset on deadline")
	}

	return nil
}
