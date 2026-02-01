package todo

type User struct {
	Id int `json:"-"`
	Name string `json:"name" binding:"required"`	// binding:"required" эти теги валидируютт поля в теле запроса и являются
	Username string `json:"username" binding:"required"` // реализацией фреймворка gin
	Password string `json:"password" binding:"required"`
}