package object

type Environment struct {
	pool  map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)

	return &Environment{pool: s, outer: nil}
}

func (env *Environment) Get(name string) (Object, bool) {
	obj, ok := env.pool[name]
	if !ok && env.outer != nil {
		obj, ok = env.outer.Get(name)
	}
	return obj, ok
}

func (env *Environment) Set(name string, val Object) Object {
	env.pool[name] = val
	return val
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
