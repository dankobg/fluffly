package server

import (
	"context"
	"net/http"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/shared"
	orykratos "github.com/ory/client-go"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListCourierMessages(ctx context.Context, request api.ListCourierMessagesRequestObject) (api.ListCourierMessagesResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "CourierMessages",
			Object:    "courier_messages",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListCourierMessages403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("message_permission", "permission denied")}, nil
	}

	req := a.Kratos.Admin.CourierAPI.ListCourierMessages(ctx)
	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}

	if request.Params.PageToken != nil && *request.Params.PageToken != "1" {
		req = req.PageToken(*request.Params.PageToken)
	}

	if request.Params.Recipient != nil {
		req = req.Recipient(*request.Params.Recipient)
	}

	if request.Params.Status != nil {
		req = req.Status(orykratos.CourierMessageStatus(*request.Params.Status))
	}

	courierMessages, courierMessagesResp, err := req.Execute()
	if err != nil {
		return api.ListCourierMessagesdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "message_list", "failed to list courier messages")}, nil
	}

	defer func() { _ = courierMessagesResp.Body.Close() }()

	resp := make(api.ListCourierMessages200JSONResponse, 0)

	for _, message := range courierMessages {
		res, err := dto.MessageToResponse(message)
		if err != nil {
			return nil, err
		}

		resp = append(resp, res)
	}

	return resp, nil
}

func (a *ApiHandler) GetCourierMessage(ctx context.Context, request api.GetCourierMessageRequestObject) (api.GetCourierMessageResponseObject, error) {
	sess := GetSession(ctx)
	req := a.Kratos.Admin.CourierAPI.GetCourierMessage(ctx, request.ID)

	courierMessage, courierMessageResp, err := req.Execute()
	if err != nil {
		return api.GetCourierMessage404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("message_not_found", "courier message not found")}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			// Namespace: "CourierMessage",
			// Object:    shared.AuthzCourierMessageID(request.ID),
			Namespace: "CourierMessages",
			Object:    "courier_messages",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetCourierMessage403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("message_permission", "permission denied")}, nil
	}

	defer func() { _ = courierMessageResp.Body.Close() }()

	resp, err := dto.MessageToResponse(*courierMessage)
	if err != nil {
		return nil, err
	}

	return api.GetCourierMessage200JSONResponse(resp), nil
}
