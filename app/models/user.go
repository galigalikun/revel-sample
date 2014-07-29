package models

import "github.com/revel/revel"

type User struct {
    Id   int
    Name string
}

func (user User) Validate(v *revel.Validation) {
    v.Required(user.Name)
}

