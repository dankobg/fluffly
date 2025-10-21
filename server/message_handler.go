package server

import (
	"context"
	"net/http"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/dto"
	"github.com/ory/client-go"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListCourierMessages(ctx context.Context, request api.ListCourierMessagesRequestObject) (api.ListCourierMessagesResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "CourierMessages",
			Object:    "courier_messages",
			Relation:  "view",
			Subject:   rts.NewSubjectID(authzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.ListCourierMessagesdefaultJSONResponse{Body: api.Error{Code: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}}, nil
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
		req = req.Status(client.CourierMessageStatus(*request.Params.Status))
	}
	courierMessages, courierMessagesResp, err := req.Execute()
	if err != nil {
		return api.ListCourierMessages400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: 400, Message: err.Error()}}, nil
	}
	defer courierMessagesResp.Body.Close()
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

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "CourierMessage",
			Object:    authzCourierMessageID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID(authzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.GetCourierMessagedefaultJSONResponse{Body: api.Error{Code: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}}, nil
	}

	req := a.Kratos.Admin.CourierAPI.GetCourierMessage(ctx, request.ID)
	courierMessage, courierMessageResp, err := req.Execute()
	if err != nil {
		return api.GetCourierMessage404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "Courier message not found"}}, nil
	}
	defer courierMessageResp.Body.Close()
	resp, err := dto.MessageToResponse(*courierMessage)
	if err != nil {
		return nil, err
	}
	return api.GetCourierMessage200JSONResponse(resp), nil
}
