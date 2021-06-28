package service

import "fmt"

func Publish(file, env, cluster string) {
	fmt.Println(env, file, cluster)
}
