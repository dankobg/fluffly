package server

import (
	"context"
	"fmt"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/shared"
	"github.com/google/uuid"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) GetAnalyticsStats(ctx context.Context, request api.GetAnalyticsStatsRequestObject) (api.GetAnalyticsStatsResponseObject, error) {
	sess := GetSession(ctx)

	stats, err := a.persistor.Analytics().GetAnalyticsStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Analytics",
			Object:    "analytics",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetAnalyticsStats403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("stats_permission", "permission denied")}, nil
	}

	resp := api.GetAnalyticsStats200JSONResponse(api.GetAnalyticsStats200JSONResponse{
		Stats: dto.AnalyticsStatsToResponse(stats),
	})

	return resp, nil
}

func (a *ApiHandler) GetMyAnalyticsStats(ctx context.Context, request api.GetMyAnalyticsStatsRequestObject) (api.GetMyAnalyticsStatsResponseObject, error) {
	sess := GetSession(ctx)

	userID := uuid.MustParse(sess.Identity.Id)
	stats, err := a.persistor.Analytics().GetMyAnalyticsStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get my stats: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Analytic",
			Object:    shared.AuthzAnalyticsID(sess.Identity.Id),
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetMyAnalyticsStats403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("my_stats_permission", "permission denied")}, nil
	}

	resp := api.GetMyAnalyticsStats200JSONResponse(api.GetMyAnalyticsStats200JSONResponse{
		Stats: dto.MyAnalyticsStatsToResponse(stats),
	})

	return resp, nil
}
