package main

import "testing"

func TestReceiptHTMLDecision(t *testing.T) {
	cases := []struct {
		name     string
		donation Donation
		want     string
	}{
		{"named donor", Donation{DonorName: "Mina", Amount: 42.5, ReceiptID: "R-7"}, "Thank you, Mina"},
		{"receipt id included", Donation{DonorName: "Lee", Amount: 10, ReceiptID: "R-8"}, "Receipt R-8"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := receiptHTML(tc.donation)
			if !contains(got, tc.want) {
				t.Fatalf("receiptHTML() = %q, want substring %q", got, tc.want)
			}
		})
	}
}

func contains(value, part string) bool {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return true
		}
	}
	return false
}
