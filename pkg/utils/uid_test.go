package utils

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestUid(t *testing.T) {
	uid, err := GenUIDForHost()

	valid := UidIsValid(uid)

	gtest.C(t, func(t *gtest.T) {
		t.Assert(err, nil)
		t.Assert(valid, true)
	})
}
