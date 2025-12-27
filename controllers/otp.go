package controllers

import (
	"context"
	"fmt"
	"service/mail-server/config"
	mail "service/mail-server/gen"
	mailservices "service/mail-server/services/mail"
	otpservices "service/mail-server/services/otp"

	"sync"
	"time"
)

type OtpService struct {
	mail.UnimplementedOtpServiceServer
	Config     *config.Config
	MailClient *mailservices.MailClient
}

type OtpUser struct {
	otp    string
	expiry time.Time
}

var OtpStore = make(map[string]OtpUser)
var StoreMu sync.RWMutex

var duration = 120

func (s *OtpService) SendCode(ctx context.Context, req *mail.SendCodeRequest) (*mail.SendCodeResponse, error) {
	usermail := req.Email
	otpCode, created := otpservices.GenerateOtp()
	expiry := created.Add(time.Duration(duration) * time.Second)

	subject := "Your Verification Code"
	body := fmt.Sprintf("Your OTP code is: %s. It expires in %d seconds.", otpCode, duration)

	err := s.MailClient.SendEmail([]string{usermail}, subject, body)
	if err != nil {
		return &mail.SendCodeResponse{
			Status: mail.Status_FAILED,
		}, nil
	}

	StoreMu.Lock()
	OtpStore[usermail] = OtpUser{
		otp:    otpCode,
		expiry: expiry,
	}
	StoreMu.Unlock()

	return &mail.SendCodeResponse{
		Status: mail.Status_SUCCESS,
	}, nil
}

func (s *OtpService) ValidateCode(ctx context.Context, req *mail.ValidateCodeRequest) (*mail.ValidateCodeResponse, error) {
	usermail := req.Email
	inputCode := req.OtpCode

	StoreMu.Lock()
	val, exists := OtpStore[usermail]
	StoreMu.Unlock()

	if !exists {
		return &mail.ValidateCodeResponse{
			Status:  mail.Status_FAILED,
			IsValid: false,
		}, nil
	}

	if time.Now().After(val.expiry) {
		delete(OtpStore, usermail)
		return &mail.ValidateCodeResponse{
			Status:  mail.Status_FAILED,
			IsValid: false,
		}, nil
	}

	isValid := val.otp == inputCode
	if isValid {
		delete(OtpStore, usermail)
	}
	return &mail.ValidateCodeResponse{
		Status:  mail.Status_SUCCESS,
		IsValid: isValid,
	}, nil
}

func StartCleanupWorker(minutes int) {
	for {
		time.Sleep(time.Duration(minutes) * time.Minute)
		start := time.Now()
		StoreMu.Lock()
		now := time.Now()
		keysDeleted := 0
		for email, user := range OtpStore {
			if now.After(user.expiry) {
				delete(OtpStore, email)
				keysDeleted++
			}
		}
		StoreMu.Unlock()
		elapsed := time.Since(start)

		if keysDeleted > 0 {
			fmt.Printf("🧹 Cleanup: Removed %d expired OTPs in %s\n", keysDeleted, elapsed)
		}
	}
}
