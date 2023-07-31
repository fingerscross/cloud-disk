package models

type UserBasic struct {
	Id       int
	Identity string
	Name     string
	Password string
	Email    string
}

func (table UserBasic) tablename() string {
	return "userbasic"
}
