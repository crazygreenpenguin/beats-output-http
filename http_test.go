// +build !integration

package http

import "testing"

func Test_maskPass(t *testing.T) {
	if maskPass("r") != "*" {
		t.Fail()
	}

	if maskPass("er") != "**" {
		t.Fail()
	}

	if maskPass("ert") != "***" {
		t.Fail()
	}

	if maskPass("hgfd") != "****" {
		t.Fail()
	}

	if maskPass("password") != "********" {
		t.Fail()
	}

	if maskPass("password1") != "pa*****d1" {
		t.Fail()
	}
}
