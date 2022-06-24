package data

type UserInterface interface {
	GetAll() ([]*User, error)
	GetByEmail(string) (*User, error)
	GetOne(int) (*User, error)
	Update(User) error
	Delete() error
	DeleteByID(int) error
	Insert(User) (int, error)
	ResetPassword(string) error
	PasswordMatches(string) (bool, error)
}

type PlanInterface interface {
	GetAll() ([]*Plan, error)
	GetOne(int) (*Plan, error)
	SubscribeUserToPlan(User, Plan) error
	AmountForDisplay() string
}
