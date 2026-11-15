package main

import "testing"

// TestVolumeTierBoundary pins ORD-04: one automatic volume adjustment of 5%
// for quantities 10-24 inclusive, 10% for quantities 25 and above, and none
// below 10. It exercises the tier-selection boundary through orderDiscount,
// the function create_order uses to price an order.
func TestVolumeTierBoundary(t *testing.T) {
	coffee := Product{SKU: "COFFEE-250", UnitPrice: 1000, Currency: "USD", MaxOrderQuantity: maxOrderQuantity}
	cases := []struct {
		quantity int
		want     int
	}{
		{quantity: 9, want: 0},
		{quantity: 10, want: 500},
		{quantity: 24, want: 1200},
		{quantity: 25, want: 2500},
	}
	for _, tc := range cases {
		got, err := orderDiscount(coffee, tc.quantity, "")
		if err != nil {
			t.Fatalf("quantity %d: unexpected error: %v", tc.quantity, err)
		}
		if got != tc.want {
			t.Errorf("quantity %d: discount_cents = %d, want %d", tc.quantity, got, tc.want)
		}
	}
}
