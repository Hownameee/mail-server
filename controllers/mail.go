package controllers

import (
	"context"
	"fmt"
	"service/mail-server/config"
	mail "service/mail-server/gen"
	mailservices "service/mail-server/services/mail"
)

type MailService struct {
	mail.UnimplementedMailServiceServer
	Config     *config.Config
	MailClient *mailservices.MailClient
}

func (s *MailService) SendMail(ctx context.Context, req *mail.SendEmailRequest) (*mail.SendEmailResponse, error) {
	if req.To == "" {
		return &mail.SendEmailResponse{
			Status: mail.Status_FAILED,
		}, fmt.Errorf("recipient email is required")
	}

	err := s.MailClient.SendEmail([]string{req.To}, req.Subject, req.Body)

	if err != nil {
		return &mail.SendEmailResponse{
			Status: mail.Status_FAILED,
		}, err
	}

	return &mail.SendEmailResponse{
		Status: mail.Status_SUCCESS,
	}, nil
}
