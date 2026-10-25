package main

import (
	"context"
	"fmt"
)

type Donation struct {
	DonorEmail string  `json:"donor_email"`
	DonorName  string  `json:"donor_name"`
	Amount     float64 `json:"amount"`
	ReceiptID  string  `json:"receipt_id"`
}

type VolunteerReminder struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Shift string `json:"shift"`
}

type CampaignReport struct {
	Campaign  string  `json:"campaign"`
	Donations int     `json:"donations"`
	Raised    float64 `json:"raised"`
}

func receiptHTML(d Donation) string {
	return fmt.Sprintf("<h1>Thank you, %s</h1><p>Receipt %s records a donation of $%.2f.</p>", d.DonorName, d.ReceiptID, d.Amount)
}

func reminderHTML(r VolunteerReminder) string {
	return fmt.Sprintf("<h1>Volunteer reminder</h1><p>Hello %s, your shift is %s.</p>", r.Name, r.Shift)
}

func reportHTML(r CampaignReport) string {
	return fmt.Sprintf("<h1>%s campaign report</h1><p>%d donations raised $%.2f.</p>", r.Campaign, r.Donations, r.Raised)
}

func sendReceipt(ctx context.Context, client *EmailClient, d Donation) (string, error) {
	return client.Send(ctx, d.DonorEmail, "Your donation receipt "+d.ReceiptID, receiptHTML(d), "receipt-"+d.ReceiptID)
}

func sendReminder(ctx context.Context, client *EmailClient, r VolunteerReminder) (string, error) {
	return client.Send(ctx, r.Email, "Volunteer shift reminder", reminderHTML(r), "reminder-"+r.Shift)
}
