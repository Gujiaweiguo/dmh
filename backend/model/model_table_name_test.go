package model

import "testing"

func TestTableNames(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"User", User{}.TableName(), "users"},
		{"Role", Role{}.TableName(), "roles"},
		{"Campaign", (&Campaign{}).TableName(), "campaigns"},
		{"Order", (&Order{}).TableName(), "orders"},
		{"Reward", (&Reward{}).TableName(), "rewards"},
		{"Distributor", Distributor{}.TableName(), "distributors"},
		{"Member", Member{}.TableName(), "members"},
		{"PosterTemplate", PosterTemplate{}.TableName(), "poster_templates"},
		{"UserFeedback", UserFeedback{}.TableName(), "user_feedback"},
		{"PasswordPolicy", PasswordPolicy{}.TableName(), "password_policies"},
		{"VerificationRecord", VerificationRecord{}.TableName(), "verification_records"},
	}

	for _, tc := range cases {
		if tc.got != tc.want {
			t.Fatalf("%s table name mismatch: got %s want %s", tc.name, tc.got, tc.want)
		}
	}
}
