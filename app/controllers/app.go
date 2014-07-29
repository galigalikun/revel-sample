package controllers

import "github.com/revel/revel"
import(
    "myapp/app/models"
    "fmt"
    "log"
)

type App struct {
    *revel.Controller
    models.Transactional
}

func (c App) Index() revel.Result {


    rows, _ := c.Txn.Select(models.User{}, "select * from user")
    for _, row := range rows {
        user := row.(*models.User)
        fmt.Printf("%d, %s\n", user.Id, user.Name)
    }
    
    return c.RenderJson(rows) // Render()
}

func (c App) Insert() revel.Result {
    user := models.User{Name: "insert"}
    user.Validate(c.Validation)
    
    if c.Validation.HasErrors() {
        c.Validation.Keep()
        c.FlashParams()
        return c.RenderJson("error")
    }
    err := c.Txn.Insert(&user)
    checkError(err, "insert error")
    return c.RenderJson(user)
}

func (c App) Update(id int) revel.Result {
    var user models.User
    err := c.Txn.SelectOne(&user, "select * from user where Id=?", id)
    checkError(err, fmt.Sprintf("update select one error:%d:", id))

    user.Name = "update"
    user.Validate(c.Validation)
    
    if c.Validation.HasErrors() {
        c.Validation.Keep()
        c.FlashParams()
        return c.RenderJson("error")
    }

    count, err := c.Txn.Update(&user)
    checkError(err, fmt.Sprintf("update error:%d", count))

    return c.RenderJson(user)
}

func (c App) Delete(id int) revel.Result {
    var user models.User
    err := c.Txn.SelectOne(&user, "select * from user where Id = ?", id)
    checkError(err, "delete error")

    count, err := c.Txn.Delete(&user)
    checkError(err, "update error")

    return c.RenderJson(count)
}

func checkError(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}
