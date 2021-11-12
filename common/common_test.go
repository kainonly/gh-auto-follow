package common

import (
	"os"
	"testing"
)

var api *API

func TestMain(m *testing.M) {
	api = NewAPI(os.Getenv("USERNAME"), os.Getenv("TOKEN"))
	os.Exit(m.Run())
}
