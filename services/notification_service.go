package services

import (
	"log"
	"time"
	"shuttle/repositories"
)

type ShuttleNotificationService struct {
	ShuttleRepo *repositories.ShuttleRepository
	FCMService  *FCMService
}

func NewShuttleNotificationService(repo *repositories.ShuttleRepository, fcmService *FCMService) *ShuttleNotificationService {
	return &ShuttleNotificationService{ShuttleRepo: repo, FCMService: fcmService}
}

func (n *ShuttleNotificationService) MonitorShuttleStatus() {
	for {
		// Monitor status "menuju sekolah" dalam 1 menit terakhir
		shuttles, err := n.ShuttleRepo.FetchShuttleTrackByParent("menuju sekolah", "1 MINUTE")
		if err != nil {
			log.Printf("Error querying shuttle status: %v", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		for _, shuttle := range shuttles {
			token := "FCM_DEVICE_TOKEN" // Ganti dengan token perangkat parent
			title := "Shuttle Status Update"
			body := shuttle.StudentName + " kini " + shuttle.ShuttleStatus

			err := n.FCMService.SendNotification(token, title, body)
			if err != nil {
				log.Printf("Failed to send notification: %v", err)
			}
		}

		time.Sleep(1 * time.Minute) // Delay sebelum iterasi berikutnya
	}
}
