package dto

import (
	"fmt"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/google/uuid"
	kratos "github.com/ory/client-go"
)

func MessageToResponse(message kratos.Message) (api.Message, error) {
	dispatches := make([]api.MessageDispatch, 0, len(message.Dispatches))
	for _, d := range message.Dispatches {
		id, err := uuid.Parse(message.Id)
		if err != nil {
			return api.Message{}, fmt.Errorf("failed to parse dispatch uuid: %w", err)
		}
		messageID, err := uuid.Parse(message.Id)
		if err != nil {
			return api.Message{}, fmt.Errorf("failed to parse dispatch messageid uuid: %w", err)
		}
		dispatches = append(dispatches, api.MessageDispatch{
			ID:        id,
			MessageID: messageID,
			Status:    api.MessageDispatchStatus(d.Status),
			Error:     &d.Error,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		})
	}
	id, err := uuid.Parse(message.Id)
	if err != nil {
		return api.Message{}, fmt.Errorf("failed to parse message uuid: %w", err)
	}
	resp := api.Message{
		ID:           id,
		Body:         message.Body,
		Subject:      message.Subject,
		Channel:      message.Channel,
		Recipient:    message.Recipient,
		Status:       api.CourierMessageStatus(message.Status),
		TemplateType: api.CourierMessageTemplateType(message.TemplateType),
		Type:         api.CourierMessageType(message.Type),
		SendCount:    message.SendCount,
		Dispatches:   &dispatches,
		CreatedAt:    message.CreatedAt,
		UpdatedAt:    message.UpdatedAt,
	}
	return resp, nil
}
