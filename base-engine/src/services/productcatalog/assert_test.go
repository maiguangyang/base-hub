package productcatalog

import "testing"

func mustCatalog[T any](t *testing.T) func(T, error) T {
	t.Helper()
	return func(value T, err error) T {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		return value
	}
}

func mustCatalogErrorFree(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
