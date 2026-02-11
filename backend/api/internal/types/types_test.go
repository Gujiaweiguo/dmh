package types

import "testing"

func TestTypeStructs(t *testing.T) {
	req := AdminGetUsersReq{Page: 1, PageSize: 10, Status: "active"}
	if req.Page != 1 || req.PageSize != 10 || req.Status != "active" {
		t.Fatalf("admin get users req invalid")
	}
	resp := CommonResp{Message: "ok"}
	if resp.Message != "ok" {
		t.Fatalf("common resp invalid")
	}
	fv := FormFieldValidation{Pattern: "^a", MinLength: 1, MaxLength: 10}
	if fv.MinLength != 1 || fv.MaxLength != 10 {
		t.Fatalf("form field validation invalid")
	}
}
