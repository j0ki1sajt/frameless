package storages

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"reflect"

	"github.com/adamluzsi/frameless/reflects"

	"github.com/adamluzsi/frameless/iterators"

	"github.com/adamluzsi/frameless"
	"github.com/adamluzsi/frameless/queryusecases"
)

func NewMemory() frameless.Storage {
	return &memory{make(map[string]memoryTable)}
}

type memory struct {
	db map[string]memoryTable
}

type memoryTable map[string]frameless.Entity

func (storage *memory) Create(e frameless.Entity) error {

	id, err := randID()

	if err != nil {
		return err
	}

	storage.tableFor(e)[id] = e
	return reflects.SetID(e, id)
}

func (storage *memory) Find(quc frameless.QueryUseCase) frameless.Iterator {
	switch quc.(type) {
	case queryusecases.ByID:
		byID := quc.(queryusecases.ByID)
		entity := storage.tableFor(byID.Type)[byID.ID]

		return iterators.NewForSingleElement(entity)

	case queryusecases.AllFor:
		byAll := quc.(queryusecases.AllFor)
		table := storage.tableFor(byAll.Type)

		entities := []frameless.Entity{}
		for _, entity := range table {
			entities = append(entities, entity)
		}

		return iterators.NewForSlice(entities)

	default:
		panic(fmt.Sprintf("%s not implemented", reflect.TypeOf(quc).Name()))

	}
}

func (storage *memory) Exec(frameless.QueryUseCase) error {
	panic("not implemented")
}

//
//
//

func (storage *memory) tableFor(e frameless.Entity) memoryTable {

	t := reflect.TypeOf(e)
	var name string

	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}

	if _, ok := storage.db[name]; !ok {
		storage.db[name] = make(memoryTable)
	}

	return storage.db[name]
}

func randID() (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

	bytes := make([]byte, 42)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}