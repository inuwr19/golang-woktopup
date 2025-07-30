package config

import (
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var MidtransClient snap.Client

func InitMidtrans() {
	MidtransClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
}
