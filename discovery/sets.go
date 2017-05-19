package discovery

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/boltdb/bolt"
	"github.com/choria-io/pdbproxy/models"
)

// Sets manages named saved queries
type Sets struct {
	DB *bolt.DB
}

// Backup performs a backup to the given path
func (s Sets) Backup(path *string) error {
	file, err := os.Create(*path)

	if err != nil {
		return err
	}

	err = s.DB.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(file)

		return err
	})

	return err
}

// Get retrieves the definition for a set
func (s Sets) Get(setName string) (*models.Set, error) {
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

// Delete removes a set from the database
func (s Sets) Delete(set string) error {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("choria"))
		b := tx.Bucket([]byte("choria"))
		err := b.Delete([]byte(set))

		return err
	})

	return err
}

// Update updates or creates a set
func (s Sets) Update(request *models.Set) error {
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

// Sets retrieve a list of known sets from the database
func (s Sets) Sets() []models.Word {
	var sets []models.Word

	db.View(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("choria"))
		b := tx.Bucket([]byte("choria"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				sets = append(sets, models.Word(k))
			}
		}

		return nil
	})

	return sets
}

// Exists determines if the set is known
func (s Sets) Exists(setName string) bool {
	set, _ := s.Get(setName)

	if set.Set != "" {
		return true
	}

	return false
}
