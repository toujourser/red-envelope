package main

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	vtzh "gopkg.in/go-playground/validator.v9/translations/zh"
)

type User struct {
	Name  string `validate:"required"`
	Age   int    `validate:"gte=0,lte=130"`
	Email string `validate:"required,email"`
}

func main() {
	user := &User{
		Name:  "tom",
		Age:   -23,
		Email: "17qq.com",
	}

	v := validator.New()
	cn := zh.New()
	uni := ut.New(cn, cn)
	translator, found := uni.GetTranslator("zh")
	if found {
		err := vtzh.RegisterDefaultTranslations(v, translator)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
	} else {
		fmt.Printf("not found")
	}
	err := v.Struct(user)
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				fmt.Println(err.Translate(translator))
			}
		}
		return
	}

}
