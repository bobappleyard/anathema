package di

import "context"

var scopeKey = &struct{}{}

func (s *Scope) Install(ctx context.Context) context.Context {
	return context.WithValue(ctx, scopeKey, s)
}

func Require(ctx context.Context, ref interface{}) error {
	scope := GetScope(ctx)
	return scope.Require(ref)
}

func GetScope(ctx context.Context) *Scope {
	return ctx.Value(scopeKey).(*Scope)
}

func EnterScope(ctx context.Context) *Scope {
	next := GetScope(ctx)
	if next == nil {
		next = BaseScope()
	}
	return &Scope{next: next}
}

func BaseScope() *Scope {
	return &Scope{rules: []Rule{
		&structRule{},
		&sliceRule{},
		&pointerRule{},
		&injectRule{},
	}}
}
