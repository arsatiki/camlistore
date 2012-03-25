/*
Copyright 2011 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mysqlindexer

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"camlistore.org/pkg/blobserver"
	"camlistore.org/pkg/index"
	"camlistore.org/pkg/jsonconfig"

	_ "camlistore.org/third_party/github.com/ziutek/mymysql/godrv"
)

type myIndexStorage struct {
	host, user, password, database string

	db *sql.DB
}

var _ index.IndexStorage = (*myIndexStorage)(nil)

func (ms *myIndexStorage) BeginBatch() index.BatchMutation {
	// TODO
	return nil 
}

func (ms *myIndexStorage) CommitBatch(b index.BatchMutation) error {
	return errors.New("TODO(bradfitz): implement")
}

func (ms *myIndexStorage) Get(key string) (value string, err error) {
	panic("TODO(bradfitz): implement")
}

func (ms *myIndexStorage) Set(key, value string) error {
	return errors.New("TODO(bradfitz): implement")
}

func (ms *myIndexStorage) Delete(key string) error {
	return errors.New("TODO(bradfitz): implement")
}

func (ms *myIndexStorage) Find(key string) index.Iterator {
	panic("TODO(bradfitz): implement")
}

func newFromConfig(ld blobserver.Loader, config jsonconfig.Obj) (blobserver.Storage, error) {
	is := &myIndexStorage{
		host:     config.OptionalString("host", "localhost"),
		user:     config.RequiredString("user"),
		password: config.OptionalString("password", ""),
		database: config.RequiredString("database"),
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// TODO: connect; set is.db
	// TODO: ping it; see if it's alive; else return err

	indexer := index.New(is)

	version, err := is.SchemaVersion()
	if err != nil {
		return nil, fmt.Errorf("error getting schema version (need to init database?): %v", err)
	}
	if version != requiredSchemaVersion {
		if os.Getenv("CAMLI_ADVERTISED_PASSWORD") != "" {
			// Good signal that we're using the dev-server script, so help out
			// the user with a more useful tip:
			return nil, fmt.Errorf("database schema version is %d; expect %d (run \"./dev-server --wipe\" to wipe both your blobs and re-populate the database schema)", version, requiredSchemaVersion)
		}
		return nil, fmt.Errorf("database schema version is %d; expect %d (need to re-init/upgrade database?)",
			version, requiredSchemaVersion)
	}

	return indexer, nil
}

func init() {
	blobserver.RegisterStorageConstructor("mysqlindexer", blobserver.StorageConstructor(newFromConfig))
}

func (mi *myIndexStorage) IsAlive() (ok bool, err error) {
	// TODO(bradfitz): something more efficient here?
	_, err = mi.SchemaVersion()
	return err == nil, err
}

func (mi *myIndexStorage) SchemaVersion() (version int, err error) {
	err = mi.db.QueryRow("SELECT value FROM meta WHERE metakey='version'").Scan(&version)
	return
}