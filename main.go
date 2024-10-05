package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"
)

type User struct {
	Name string `json:"name" binding:"required,max=3,gt=0"`
	Age  int    `json:"age" binding:"gt=1"`
}

func main() {
	r := gin.Default()
	r.POST("/add", add)
	panic(r.Run(":11000"))
}

func add(c *gin.Context) {
	var u User

	if err := validErr(c.ShouldBindJSON(&u)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "ok"})
}

func validErr(err error) error {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}
	zh2 := zh.New()
	uni := ut.New(zh2, zh2)
	trans, ok := uni.GetTranslator("zh")
	if !ok {
		fmt.Printf("%s not found", "zh")
	}
	v, ok := binding.Validator.Engine().(*validator.Validate)
	zhTrans.RegisterDefaultTranslations(v, trans)

	var returnErr error
	for _, err := range errs {
		if returnErr == nil {
			returnErr = errors.New(err.Translate(trans))
		}
		fmt.Printf("%v\n", err)
	}
	return returnErr
}
