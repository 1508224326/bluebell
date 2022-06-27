package main

import (
	"fmt"
	"github.com/aiyouyo/bluebell/pkg/jwt"
)

func main() {
	token, _ := jwt.GenToken("张三", 23)
	fmt.Println(token)

	c, _ := jwt.ParseToken(token)

	fmt.Println(c)

}
