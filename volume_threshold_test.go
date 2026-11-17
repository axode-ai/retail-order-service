package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestVolumeThresholdAcceptedOutcome(t *testing.T) {
	// ORD-04 and ORD-06: explicit whole-order cents at both inclusive
	// thresholds and their immediate neighbors; no rounding is involved.
	for _, tc := range []struct {
		quantity, discount, total int
	}{
		{9, 0, 9000},
		{10, 500, 9500},
		{11, 550, 10450},
		{24, 1200, 22800},
		{25, 2500, 22500},
		{26, 2600, 23400},
	} {
		for _, code := range []string{"", "WELCOME100"} {
			t.Run(fmt.Sprintf("quantity_%d/promotion_%s", tc.quantity, code), func(t *testing.T) {
				s, store := testService()
				request := CheckoutRequest{OrderID: "volume-boundary", SKU: "COFFEE-250", Quantity: tc.quantity, PromotionCode: code}
				wantDiscount, wantTotal := tc.discount, tc.total
				// ORD-05: the fixed 100-cent promotion follows volume pricing.
				if code == "WELCOME100" {
					wantDiscount += 100
					wantTotal -= 100
				}
				response, receipt := submit(t, s, request)
				if response.Code != http.StatusCreated {
					t.Fatalf("acceptance status=%d, want 201: %+v", response.Code, receipt)
				}
				stored, ok := store.orders[request.OrderID]
				if !ok {
					t.Fatal("accepted order missing from store")
				}
				if stored.Receipt.Discount != wantDiscount || stored.Ledger.Total != wantTotal {
					t.Errorf("persisted volume outcome: discount=%d total=%d, want discount=%d total=%d", stored.Receipt.Discount, stored.Ledger.Total, wantDiscount, wantTotal)
				}
				if stored.Request != request || stored.Receipt.UnitPrice != 1000 || stored.Receipt.Currency != "USD" || stored.Ledger.Writes != 1 {
					t.Errorf("unexpected accepted snapshot: %+v", stored)
				}
				if receipt.Discount != wantDiscount || receipt.LedgerTotal == nil || *receipt.LedgerTotal != wantTotal || receipt.LedgerWrites == nil || *receipt.LedgerWrites != 1 {
					t.Errorf("response volume outcome: %+v; want discount=%d total=%d writes=1", receipt, wantDiscount, wantTotal)
				}
			})
		}
	}
}
