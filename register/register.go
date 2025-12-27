package register

import (
	"service/mail-server/config"
	"service/mail-server/controllers"
	mail "service/mail-server/gen"
	mailservices "service/mail-server/services/mail"

	"google.golang.org/grpc"
)

func Register(server *grpc.Server, config *config.Config) error {
	mailClient := mailservices.NewMailClient(config, 5)

	mail.RegisterOtpServiceServer(server, &controllers.OtpService{Config: config, MailClient: mailClient})
	go controllers.StartCleanupWorker(5)

	mail.RegisterMailServiceServer(server, &controllers.MailService{Config: config, MailClient: mailClient})
	return nil
}
