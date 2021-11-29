package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"tianxu.xin/phone/cloud"
)

func main() {
	account := os.Getenv("ICLOUD_ACCOUNT")
	password := os.Getenv("ICLOUD_PASSWORD")
	ctx := context.Background()
	client := cloud.DefaultICloud(account, password)
	if err := client.Login(ctx); err != nil {
		fmt.Println("login err", err)
	}
	photos, err := client.PhotoService(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	infos, err := photos.All(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("infos len", len(infos))

	for _, v := range infos {
		if err = v.WriteTo(ctx, v.Name()); err != nil {
			log.Fatalln(err)
		}
	}
}
