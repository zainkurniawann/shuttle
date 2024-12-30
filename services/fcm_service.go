package services

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type FCMService struct {
	App *firebase.App
}

// NewFCMService membuat instance baru dari FCMService
func NewFCMService(serviceAccountPath string) (*FCMService, error) {
	opt := option.WithCredentialsFile(serviceAccountPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	return &FCMService{App: app}, nil
}

// SendNotification mengirimkan notifikasi push ke device menggunakan token FCM
func (f *FCMService) SendNotification(token, title, body string) error {
	client, err := f.App.Messaging(context.Background())
	if err != nil {
		return err
	}

	// Gunakan tipe messaging.Message dan messaging.Notification
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token,
	}

	// Kirim pesan ke Firebase Cloud Messaging
	response, err := client.Send(context.Background(), message)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return err
	}
	log.Printf("Notification sent successfully: %s", response)
	return nil
}
