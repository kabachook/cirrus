package bbolt

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/kabachook/cirrus/pkg/database"
	"github.com/kabachook/cirrus/pkg/provider"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

var bucketSnapshot = []byte("snapshots")

type Database struct {
	database.Database
	filename string
	options  *bolt.Options
	db       *bolt.DB
	logger   *zap.Logger
}

type Config struct {
	Filename string
	Options  bolt.Options
	Logger   *zap.Logger
}

func New(cfg Config) *Database {
	return &Database{
		filename: cfg.Filename,
		options:  &cfg.Options,
		logger:   cfg.Logger,
	}
}

func (D *Database) Open() error {
	db, err := bolt.Open(D.filename, 0600, D.options)
	if err != nil {
		return err
	}
	D.db = db
	D.logger.Info("Database opened")

	err = D.db.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucketIfNotExists(bucketSnapshot)
		return err
	})
	return err
}

func (D *Database) Close() error {
	return D.db.Close()
}

func timestampToBytes(timestamp int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, timestamp)
	return buf[:n]
}

func bytesToTimestamp(raw []byte) int64 {
	// Please be safe
	t, _ := binary.Varint(raw)
	return t
}

func (D *Database) Store(timestamp int64, endpoints []provider.Endpoint) error {
	D.logger.Debug("Storing snapshot", zap.Int64("timestamp", timestamp), zap.Any("endpoints", endpoints))
	err := D.db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(bucketSnapshot)
		endpointsBytes, err := json.Marshal(endpoints)
		if err != nil {
			return err
		}
		err = b.Put(timestampToBytes(timestamp), endpointsBytes)
		return err
	})
	return err
}

func (D *Database) List() ([]database.Snapshot, error) {
	var snapshots = make([]database.Snapshot, 0)
	err := D.db.View(func(t *bolt.Tx) error {
		b := t.Bucket(bucketSnapshot)
		err := b.ForEach(func(k, v []byte) error {
			var endpoints []provider.Endpoint
			ts := bytesToTimestamp(k)
			D.logger.Debug("Get", zap.Int64("ts", ts))
			err := json.Unmarshal(v, &endpoints)
			if err != nil {
				return err
			}
			snapshots = append(snapshots, database.Snapshot{
				Timestamp: ts,
				Endpoints: endpoints,
			})
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return snapshots, nil
}

func (D *Database) Get(timestamp int64) (*database.Snapshot, error) {
	var endpoints []provider.Endpoint
	err := D.db.View(func(t *bolt.Tx) error {
		b := t.Bucket(bucketSnapshot)
		bb := b.Get(timestampToBytes(timestamp))
		if bb == nil {
			return fmt.Errorf("snapshot not found")
		}
		json.Unmarshal(bb, &endpoints)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &database.Snapshot{
		Timestamp: timestamp,
		Endpoints: endpoints,
	}, nil
}
