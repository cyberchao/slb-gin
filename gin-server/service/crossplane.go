package service

import (
	"encoding/json"
	"fmt"
	"github.com/aluttik/go-crossplane"
)

func CrossPlane() {
	payload, err := crossplane.Parse("tmp/data.txt", &crossplane.ParseOptions{})
	if payload.Status == "ok" {
		fmt.Println(payload.Config[0].Parsed[0].Block)
	}

	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
