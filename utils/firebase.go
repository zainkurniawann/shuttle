package utils

import (
	"context"
	"log"

	"firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebase() {
	opt := option.WithCredentialsFile("./firebase_key/shuttle-7e595-firebase-adminsdk-71s8d-a8fc854794.json") // Ganti dengan path ke file JSON service account
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}
	FirebaseApp = app
}