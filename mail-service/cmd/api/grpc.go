package main

import (
	"context"
	"github.com/arogyaGurkha/gurkhaland-proto/mail-service/mail"
	"log"
)

type MailServer struct {
	mail.UnimplementedMailServiceServer
	Mailer Mail
}

func (m *MailServer) SendMail(ctx context.Context, req *mail.MailRequest) (*mail.MailResponse, error) {
	input := req.GetMailReq()

	msg := Message{
		From:    input.From,
		To:      input.To,
		Subject: input.Subject,
		Data:    input.Message,
	}

	err := m.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res := &mail.MailResponse{
		Result: "Mail sent to " + input.To,
	}
	return res, nil
}
