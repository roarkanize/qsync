package qsync

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"net/url"
	"os"
)

type DB struct{ *bolt.DB }
type Tx struct{ *bolt.Tx }

func InitDB(path string) (*DB, error) {
	db, err := bolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("entries"))
		if err != nil {
			return fmt.Errorf("create `entries` bucket: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte("remotes"))
		if err != nil {
			return fmt.Errorf("create `remotes` bucket: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte("qubes"))
		if err != nil {
			return fmt.Errorf("create `qubes` bucket: %s", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (tx *Tx) AddEntry(name []byte, path []byte, remote []byte) error {
	// note: maybe sanitize the entry key? something like slugify or whatever.
	err := tx.Bucket([]byte("entries")).Put(name, path)
	if err != nil {
		return err
	}

	_, err = url.ParseRequestURI(string(remote))
	if err != nil {
		return errors.New("remote string is not valid URI")
	}

	err = tx.Bucket([]byte("remotes")).Put(name, remote)
	if err != nil {
		return err
	}

	return nil
}

func (tx *Tx) AddQube(name []byte) error {
	// check if VM exists
	exists, err := VMExists(name)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("qube '%s' does not exist", name)
	}

	// note: is there a better way to get current VM name?
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	if bytes.Compare(name, []byte(hostname)) == 0 {
		return errors.New("cannot add sync qube to database")
	}

	err = tx.Bucket([]byte("qubes")).Put(name, nil)
	if err != nil {
		return err
	}

	return nil
}
