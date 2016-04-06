package datastore

import (
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Entity interface {
	EntityName() string
	SetKey(key *datastore.Key)
	GetKey() *datastore.Key
	GetID() int64
}

////////////

type Keyable struct {
	ID  int64          `datastore:"-" json:"id"`
	Key *datastore.Key `datastore:"-" json:"-"`
}

func (x *Keyable) GetID() int64 {
	return x.ID
}

func (x *Keyable) SetKey(key *datastore.Key) {
	x.Key = key
	x.ID = key.IntID()
}

func (x *Keyable) GetKey() *datastore.Key {
	return x.Key
}

////////////

func Find(ctx context.Context, dst Entity, sid string) error {
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return err
	}

	key := datastore.NewKey(ctx, dst.EntityName(), "", id, nil)
	err = datastore.Get(ctx, key, dst)
	if err != nil {
		return err
	}
	dst.SetKey(key)
	return nil
}

func Create(ctx context.Context, dst Entity) error {
	key, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, dst.EntityName(), nil), dst)
	if err != nil {
		return err
	}
	dst.SetKey(key)
	return nil
}

func Update(ctx context.Context, dst Entity) error {
	key, err := datastore.Put(ctx, dst.GetKey(), dst)
	if err != nil {
		return err
	}
	dst.SetKey(key)
	return nil
}

func Delete(ctx context.Context, dst Entity) error {
	return datastore.Delete(ctx, dst.GetKey())
}
