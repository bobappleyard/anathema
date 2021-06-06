package di

import (
	"context"
	"reflect"
	"sync"
)

var (
	contextType = reflect.TypeOf(new(context.Context)).Elem()
	lockMark    = new(struct{})
)

type factoryStrategy struct {
	repo *factoryRepository
}

type factoryRepository struct {
	lock      sync.RWMutex
	factories []factory
}

type factory struct {
	forType reflect.Type
	impl    func(ctx context.Context) (reflect.Value, error)
	cache   reflect.Value
}

func (s *factoryStrategy) furnishValue(ctx context.Context, v reflect.Value) (bool, error) {
	fs := s.repo.findFactories(ctx, v.Type())

	switch len(fs) {
	case 0:
		return false, nil

	case 1:
		collab, err := s.repo.triggerFactory(ctx, fs[0])
		if err != nil {
			return true, err
		}

		v.Set(collab)
		return true, nil

	default:
		return true, ErrTooManyFactories
	}
}

func (r *factoryRepository) findFactories(ctx context.Context, t reflect.Type) []*factory {
	if !r.inLock(ctx) {
		r.lock.RLock()
		defer r.lock.RUnlock()
	}

	var fs []*factory
	for i, f := range r.factories {
		if f.forType.AssignableTo(t) {
			fs = append(fs, &r.factories[i])
		}
	}

	return fs
}

func (r *factoryRepository) triggerFactory(ctx context.Context, f *factory) (reflect.Value, error) {
	v := r.checkFactoryCache(ctx, f)
	if v.IsValid() {
		return v, nil
	}

	if !r.inLock(ctx) {
		r.lock.Lock()
		defer r.lock.Unlock()
	}

	if !f.cache.IsValid() {
		collab, err := f.impl(r.enterLock(ctx))
		if err != nil {
			return reflect.Value{}, err
		}
		f.cache = collab
	}

	return f.cache, nil
}

func (r *factoryRepository) checkFactoryCache(ctx context.Context, f *factory) reflect.Value {
	if !r.inLock(ctx) {
		r.lock.RLock()
		defer r.lock.RUnlock()
	}

	return f.cache
}

func (r *factoryRepository) inLock(ctx context.Context) bool {
	return ctx.Value(lockMark) == r
}

func (r *factoryRepository) enterLock(ctx context.Context) context.Context {
	return context.WithValue(ctx, lockMark, r)
}
