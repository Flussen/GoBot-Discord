package commands

import "testing"

func TestPasswordGenerator(t *testing.T) {
	leng := 10
	value := PasswordGenerator(leng)

	if len(value) != leng {
		t.Errorf("password is below to %d\n, password: %s", leng, value)
	}
	t.Logf("\nPassword: %s\nLength: %d", value, len(value))
}
