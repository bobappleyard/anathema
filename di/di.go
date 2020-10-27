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
	return &Scope{next: GetScope(ctx)}
}

func BaseScope() *Scope {
	return &Scope{providers: []Provider{
		&pointerProvider{},
		&sliceProvider{},
		&structProvider{},
	}}
}
