// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bookshelf

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func testDB(t *testing.T, db BookDatabase) {
	t.Helper()

	defer db.Close()

	b := &Book{
		Author:      "testy mc testface",
		Title:       fmt.Sprintf("t-%d", time.Now().Unix()),
		Description: "desc",
	}

	id, err := db.AddBook(b)
	if err != nil {
		t.Fatal(err)
	}

	b.ID = id
	b.Description = "newdesc"
	if err := db.UpdateBook(b); err != nil {
		t.Error(err)
	}

	gotBook, err := db.GetBook(id)
	if err != nil {
		t.Error(err)
	}
	if got, want := gotBook.Description, b.Description; got != want {
		t.Errorf("Update description: got %q, want %q", got, want)
	}

	if err := db.DeleteBook(id); err != nil {
		t.Error(err)
	}

	if _, err := db.GetBook(id); err == nil {
		t.Error("want non-nil err")
	}
}

func TestMemoryDB(t *testing.T) {
	testDB(t, newMemoryDB())
}

func TestDatastoreDB(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	db, err := newDatastoreDB(client)
	if err != nil {
		t.Fatal(err)
	}
	testDB(t, db)
}
