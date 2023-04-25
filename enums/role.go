package enums

const (
	User       string = "user"
	Admin      string = "admin"
	SuperAdmin string = "superadmin"
)

type Roles struct{}

func (r Roles) User() string {
	return User
}

func (r Roles) Admin() string {
	return Admin
}

func (r Roles) SuperAdmin() string {
	return SuperAdmin
}
