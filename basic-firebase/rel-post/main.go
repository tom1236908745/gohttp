package rel_post

import (
	"net/http"

	"github.com/labstack/echo"
)

type User struct {
	Name    string `json:"name"`
	Age     string `json:"age"`
	Address string `json:"address"`
}

func main() {
	e := echo.New()
	e.POST("/users", addUser)
	e.Logger.Fatal(e.Start(":1323"))
}

func addUser(c echo.Context) error {
	u := new(User)
	if error := c.Bind(u); error != nil {
		return error
	}
	// データ追加呼び出し
	if error := dataAdd(u.Name, u.Age, u.Address); error != nil {
		return error
	}

	return c.JSON(http.StatusOK, u)
}
