package controllers

import (
	"log"
	"membership-fitness-centre/services"
	"membership-fitness-centre/utils"
)

type SubscriptionController struct {
	subscriptionService *services.Subscriptionervice
}

func NewSubscriptionController(serviceSubscription *services.Subscriptionervice) *SubscriptionController {
	return &SubscriptionController{subscriptionService: serviceSubscription}
}

func (c *SubscriptionController) CheckSubscriptions() {
	members, err := c.subscriptionService.GetExpiringSubscriptions()
	if err != nil {
		log.Fatal(err)
	}

	for _, member := range members {
		msg := []byte("Subject: Subscription Expiry Alert!\n\nYour subscription is expiring soon!")
		err := utils.SendPackageEmails(member.Email, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}
