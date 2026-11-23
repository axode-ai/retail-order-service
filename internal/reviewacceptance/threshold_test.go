package reviewacceptance

import "testing"

func TestMeetsMinimum(t *testing.T) {
	if MeetsMinimum(9, 10) {
		t.Error("amount below the minimum must fail")
	}
	if !MeetsMinimum(10, 10) {
		t.Error("amount equal to the minimum must qualify")
	}
	if !MeetsMinimum(11, 10) {
		t.Error("amount above the minimum must qualify")
	}
}
