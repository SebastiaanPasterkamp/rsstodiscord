package checker_test

import (
	"reflect"
	"testing"

	"github.com/SebastiaanPasterkamp/rsstodiscord/internal/checker"
)

func TestTranslate(t *testing.T) {
	t.Parallel()

	m := checker.Translate(&testItem)

	if !reflect.DeepEqual(m, testMessage) {
		t.Errorf("Incorrect message. Expected %v, got %v.", testMessage, m)
	}
}
