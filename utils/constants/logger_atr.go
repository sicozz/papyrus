package constants

type Layer string

const (
	Main       Layer = "MAIN"
	Delivery   Layer = "DELIVERY"
	Usecase    Layer = "USECASE"
	Repository Layer = "REPOSITORY"
	Utils      Layer = "UTILS"
)

type Domain string

const (
	None      Domain = ""
	User      Domain = "USER"
	Role      Domain = "ROLE"
	UserState Domain = "USER_STATE"
)
