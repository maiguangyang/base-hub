package customer

import "testing"

func TestNormalizeCNPhone(t *testing.T) {
	for input, want := range map[string]string{
		"13800000001":     "+8613800000001",
		"138-0000 0001":   "+8613800000001",
		"+86 13800000001": "+8613800000001",
		"8613800000001":   "+8613800000001",
	} {
		got, err := NormalizeCNPhone(input)
		if err != nil || got != want {
			t.Fatalf("NormalizeCNPhone(%q) = %q, %v", input, got, err)
		}
	}
	for _, input := range []string{"", "12800000001", "1380000000", "+8213800000001", "138x00000001", "138000000011"} {
		if _, err := NormalizeCNPhone(input); err == nil {
			t.Fatalf("invalid phone %q was accepted", input)
		}
	}
}

func TestMaskCNPhone(t *testing.T) {
	if got := MaskCNPhone("+8613800000001"); got != "+86138****0001" {
		t.Fatalf("masked phone = %q", got)
	}
}

func TestNormalizeCNPhonePrefix(t *testing.T) {
	if got, err := NormalizeCNPhonePrefix("138-0000"); err != nil || got != "+861380000" {
		t.Fatalf("phone prefix = %q, %v", got, err)
	}
	if got, err := NormalizeCNPhonePrefix("861380000"); err != nil || got != "+861380000" {
		t.Fatalf("country-prefixed phone = %q, %v", got, err)
	}
	for _, input := range []string{"138000", "+86138000", "2380000", "138000000001"} {
		if _, err := NormalizeCNPhonePrefix(input); err == nil {
			t.Fatalf("invalid prefix %q was accepted", input)
		}
	}
}
