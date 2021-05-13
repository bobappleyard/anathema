package di

type builder struct {
	constructors []Constructor
	fallback     []Constructor
	mutators     []Mutator
	complete     bool
	cache        []string
}

func (t *builder) Constructor(c Constructor) {
	t.constructors = append(t.constructors, c)
}

func (t *builder) FallbackConstructor(c Constructor) {
	t.fallback = append(t.fallback, c)
}

func (t *builder) Mutator(m Mutator) {
	t.mutators = append(t.mutators, m)
}

func (t *builder) Complete() {
	t.complete = true
}

func (t *builder) Cache(scope string) {
	t.cache = append(t.cache, scope)
}

func (t *builder) getConstructors() []Constructor {
	if t.constructors != nil {
		return t.constructors
	}
	return t.fallback
}
