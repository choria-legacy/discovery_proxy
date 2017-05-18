package discovery

import (
	"encoding/json"
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/choria-io/pdbproxy/models"
)

type Sets struct {
	DB *bolt.DB
}

func (s Sets) GetSet(setName string) (*models.Set, error) {
	log.Infof("Retrieving set %s", setName)

	set := models.Set{}

	err := s.DB.View(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("choria"))
		b := tx.Bucket([]byte("choria"))
		val := b.Get([]byte(setName))

		if len(val) == 0 {
			return errors.New("Could not find set " + setName)
		}

		json.Unmarshal(val, &set)

		return nil
	})

	return &set, err
}

func (s Sets) Delete(set string) error {
	log.Info("Deleting set %s", set)

	err := s.DB.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("choria"))
		b := tx.Bucket([]byte("choria"))
		err := b.Delete([]byte(set))

		return err
	})

	return err
}

func (s Sets) Update(request *models.Set) error {
	log.Infof("Updating set %s", request.Set)

	err := s.DB.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("choria"))
		b := tx.Bucket([]byte("choria"))

		j, err := json.Marshal(request)
		if err != nil {
			return err
		}

		err = b.Put([]byte(request.Set), j)

		return err
	})

	return err
}
