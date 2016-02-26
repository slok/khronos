package storage

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

// tearDownBoltDB closes and deletes boltdb
func tearDownBoltDB(db *bolt.DB) error {
	p := db.Path()
	err := db.Close()
	if err != nil {
		return err
	}

	return os.Remove(p)
}

func randomPath() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("/tmp/khronos_boltdb_test_%d.db", r.Int())
}

func TestBoltDBConnection(t *testing.T) {
	boltPath := randomPath()
	fmt.Println(boltPath)
	// Create a new boltdb connection
	c, err := NewBoltDB(boltPath, 2*time.Second)

	if err != nil {
		t.Errorf("Error creating bolt connection: %v", err)
	}

	// Check root buckets are present
	checkBuckets := []string{jobsBucket, resultsBucket}
	err = c.DB.View(func(tx *bolt.Tx) error {
		for _, cb := range checkBuckets {
			if b := tx.Bucket([]byte(cb)); b == nil {
				t.Errorf("Bucket %s not present", cb)
			}
		}
		return nil
	})

	if err != nil {
		t.Error(err)
	}

	// Close ok
	if err := tearDownBoltDB(c.DB); err != nil {
		t.Error(err)
	}

}
