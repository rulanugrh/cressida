package service

import (
	"context"
	"fmt"

	"github.com/rulanugrh/cressida/internal/entity/web"
	"github.com/rulanugrh/cressida/internal/helper"
	"github.com/rulanugrh/cressida/internal/repository"
)

type NotificationService interface {
	GetNotificationByUserID(userID uint) (*[]web.Notification, error)
}

type notification struct {
	repository repository.NotificationRepository
	log        helper.ILog
	trace      helper.IOpenTelemetry
}

func NewNotificationService(repository repository.NotificationRepository) NotificationService {
	return &notification{
		repository: repository,
		log:        helper.NewLogger(),
		trace:      helper.NewOpenTelemetry(),
	}
}

func (n *notification) GetNotificationByUserID(userID uint) (*[]web.Notification, error) {
	// span for tracing request this endpoint
	span := n.trace.StartTracer(context.Background(), "GetNotificationByUserID")
	defer span.End()

	// parsing userID for get data
	data, err := n.repository.GetNotificationByUserId(int(userID))
	if err != nil {
		n.log.Error(fmt.Sprintf("[SERVICE] - [GetNotificationByUserID] failure get nofiticatoin by this userID: %d", userID))
		return nil, web.BadRequest(err.Error())
	}

	var responses []web.Notification
	for _, v := range *data {
		responses = append(responses, web.Notification{
			Content: v.Content,
			OrderID: v.OrderID,
			Status:  v.Status,
		})
	}

	return &responses, nil
}
