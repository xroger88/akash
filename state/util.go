package state

import (
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/ovrclk/akash/types"
	"github.com/tendermint/iavl"
	tmdb "github.com/tendermint/tmlibs/db"
)

const (
	dbCacheSize = 256
	dbBackend   = tmdb.GoLevelDBBackend
)

func LoadDB(pathname string) (DB, error) {

	if pathname == "" {
		return NewMemDB(), nil
	}

	pathname, err := filepath.Abs(pathname)
	if err != nil {
		return nil, err
	}

	pathname = strings.TrimSuffix(pathname, path.Ext(pathname))

	dir := path.Dir(pathname)
	name := path.Base(pathname)

	db := tmdb.NewDB(name, dbBackend, dir)
	tree := iavl.NewVersionedTree(db, dbCacheSize)
	if _, err := tree.Load(); err != nil {
		return nil, err
	}

	mtx := new(sync.RWMutex)
	return &iavlDB{tree, mtx}, nil
}

func LoadState(db DB, gen *types.Genesis) (State, error) {

	state := NewState(db)

	if gen == nil || !db.IsEmpty() {
		return state, nil
	}

	accounts := state.Account()

	for idx := range gen.Accounts {
		if err := accounts.Save(&gen.Accounts[idx]); err != nil {
			return nil, err
		}
	}

	return state, nil
}

func NewMemDB() DB {
	mtx := new(sync.RWMutex)
	tree := iavl.NewVersionedTree(tmdb.NewMemDB(), 0)
	return &iavlDB{tree, mtx}
}

func saveObject(db DB, key []byte, obj proto.Message) error {
	buf, err := proto.Marshal(obj)
	if err != nil {
		return err
	}

	db.Set(key, buf)
	return nil
}
