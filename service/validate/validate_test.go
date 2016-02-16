package validate

import "testing"

func TestValidCron(t *testing.T) {

	tests := []struct {
		givenCron string
		wantError bool
	}{
		{givenCron: "@daily", wantError: false},
		{givenCron: "@every 1m", wantError: false},
		{givenCron: "* * * * * ?", wantError: false},
		{givenCron: "0 30 * * * *", wantError: false},
		{givenCron: "* 1", wantError: true},
		{givenCron: "every 1m", wantError: true},
		{givenCron: "* * * * * .", wantError: true},
	}

	for _, test := range tests {
		err := ValidCron(test.givenCron)

		if test.wantError && err == nil {
			t.Errorf("'%s' should raise error", test.givenCron)
		}

		if !test.wantError && err != nil {
			t.Errorf("'%s' should not raise error", test.givenCron)
		}
	}
}
