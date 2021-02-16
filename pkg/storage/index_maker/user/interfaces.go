package user

type logger interface {
	Warn(args ...interface{})
}
