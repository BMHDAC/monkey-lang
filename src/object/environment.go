package object

type Environment struct {
	pool map[string]Object
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)

	return &Environment{pool: s}
}

func (env *Environment) Get(name string) (Object, bool) {
	obj, ok := env.pool[name]
	return obj, ok
}

func (env *Environment) Set(name string, val Object) Object {
	env.pool[name] = val
	return val
}
