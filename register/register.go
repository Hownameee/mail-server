package register

import (
	"service/mail-server/config"
	"service/mail-server/controllers"
	mail "service/mail-server/gen"
	mailservices "service/mail-server/services/mail"

	"google.golang.org/grpc"
)

func Register(server *grpc.Server, cfg *config.Config) error {
	mailClient := mailservices.NewMailClient(cfg, cfg.Workers)

	mail.RegisterOtpServiceServer(server, &controllers.OtpService{Config: cfg, MailClient: mailClient})
	mail.RegisterMailServiceServer(server, &controllers.MailService{Config: cfg, MailClient: mailClient})

	go controllers.StartCleanupWorker(cfg.OtpCleanupMinutes)

	return nil
}
