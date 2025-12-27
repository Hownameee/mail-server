package mailservices

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"service/mail-server/config"
	"time"
)

type job struct {
	To      []string
	Subject string
	Body    string
	err     chan error
}

type MailClient struct {
	jobs    chan job
	cfg     *config.Config
	workers int
}

func NewMailClient(cfg *config.Config, workers int) *MailClient {
	client := &MailClient{
		jobs:    make(chan job, 100),
		cfg:     cfg,
		workers: workers,
	}

	for i := 0; i < client.workers; i++ {
		go client.worker(i + 1)
	}

	return client
}

func (m *MailClient) SendEmail(to []string, title string, content string) error {
	err := make(chan error)
	m.jobs <- job{
		To:      to,
		Subject: title,
		Body:    content,
		err:     err,
	}
	res := <-err
	return res
}

func (m *MailClient) worker(id int) {
	var client *smtp.Client

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	reconnect := func() error {
		if client != nil {
			client.Quit()
			client = nil
		}

		addr := fmt.Sprintf("%s:%s", m.cfg.EmailHost, m.cfg.EmailPort)
		conn, err := tls.Dial("tcp", addr, &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         m.cfg.EmailHost,
		})
		if err != nil {
			return err
		}

		c, err := smtp.NewClient(conn, m.cfg.EmailHost)
		if err != nil {
			return err
		}

		auth := smtp.PlainAuth("", m.cfg.EmailUser, m.cfg.EmailPass, m.cfg.EmailHost)
		if err = c.Auth(auth); err != nil {
			c.Close()
			return err
		}

		client = c
		log.Printf("[Worker %d] Connected to %s", id, m.cfg.EmailHost)
		return nil
	}

	if err := reconnect(); err != nil {
		log.Printf("[Worker %d] Initial connection failed (will retry): %v", id, err)
	}

	for {
		select {

		case j := <-m.jobs:
			if client == nil {
				if err := reconnect(); err != nil {
					log.Printf("[Worker %d] Failed to reconnect: %v. Dropping email.", id, err)

					j.err <- err
					continue
				}
			}

			err := func() error {
				if err := client.Mail(m.cfg.EmailUser); err != nil {
					return err
				}
				for _, to := range j.To {
					if err := client.Rcpt(to); err != nil {
						return err
					}
				}
				w, err := client.Data()
				if err != nil {
					return err
				}

				msg := "To: " + j.To[0] + "\r\n" +
					"Subject: " + j.Subject + "\r\n" +
					"\r\n" +
					j.Body + "\r\n"

				if _, err = w.Write([]byte(msg)); err != nil {
					return err
				}
				return w.Close()
			}()

			if err != nil {
				log.Printf("[Worker %d] Error sending email: %v. Reconnecting...", id, err)
				client = nil

				j.err <- err
			} else {
				// log.Printf("[Worker %d] Email sent to %s", id, j.To[0])

				j.err <- nil
			}

		case <-ticker.C:
			if client != nil {
				if err := client.Noop(); err != nil {
					log.Printf("[Worker %d] Heartbeat failed: %v. Connection dead.", id, err)
					client = nil
				} else {
					// log.Printf("[Worker %d] Heartbeat success!", id)
				}
			}
		}
	}
}
