package cloud

import (
	"go.uber.org/zap"
	"tianxu.xin/phone/cloud/pkg"
)

func DefaultICloud(account, password string) pkg.Client {
	log := zap.L()
	client, err := pkg.NewICloud(account, password, "./", log)
	if err != nil {
		log.Fatal("init icloud failed", zap.Error(err))
	}
	return client
}
