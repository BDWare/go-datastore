// Copyright for portions of this fork are held by [Juan Batiz-Benet, 2016] as
// part of the original go-datastore project. All other copyright for
// this fork are held by [The BDWare Authors, 2020]. All rights reserved.
// Use of this source code is governed by MIT license that can be
// found in the LICENSE file.

package retrystore

import (
	"fmt"
	"strings"
	"testing"

	failstore "github.com/bdware/go-datastore/failstore"
	key "github.com/bdware/go-datastore/key"
	dstest "github.com/bdware/go-datastore/test"
)

func TestRetryFailure(t *testing.T) {
	myErr := fmt.Errorf("this is an actual error")
	var count int
	fstore := failstore.NewFailstore(dstest.NewMapDatastoreForTest(t), func(op string) error {
		count++
		return myErr
	})

	rds := &Datastore{
		Batching: fstore,
		Retries:  5,
		TempErrFunc: func(err error) bool {
			return err == myErr
		},
	}

	k := key.NewStrKey("test")

	_, err := rds.Get(k)
	if err == nil {
		t.Fatal("expected this to fail")
	}

	if !strings.Contains(err.Error(), "ran out of retries") {
		t.Fatal("got different error than expected: ", err)
	}

	if count != 6 {
		t.Fatal("expected five retries (six executions), got: ", count)
	}
}

func TestRealErrorGetsThrough(t *testing.T) {
	myErr := fmt.Errorf("this is an actual error")
	fstore := failstore.NewFailstore(dstest.NewMapDatastoreForTest(t), func(op string) error {
		return myErr
	})

	rds := &Datastore{
		Batching: fstore,
		Retries:  5,
		TempErrFunc: func(err error) bool {
			return false
		},
	}

	k := key.NewStrKey("test")
	_, err := rds.Get(k)
	if err != myErr {
		t.Fatal("expected my own error")
	}

	_, err = rds.Has(k)
	if err != myErr {
		t.Fatal("expected my own error")
	}

	err = rds.Put(k, nil)
	if err != myErr {
		t.Fatal("expected my own error")
	}
}

func TestRealErrorAfterTemp(t *testing.T) {
	myErr := fmt.Errorf("this is an actual error")
	tempErr := fmt.Errorf("this is a temp error")
	var count int
	fstore := failstore.NewFailstore(dstest.NewMapDatastoreForTest(t), func(op string) error {
		count++
		if count < 3 {
			return tempErr
		}

		return myErr
	})

	rds := &Datastore{
		Batching: fstore,
		Retries:  5,
		TempErrFunc: func(err error) bool {
			return err == tempErr
		},
	}

	k := key.NewStrKey("test")
	_, err := rds.Get(k)
	if err != myErr {
		t.Fatal("expected my own error")
	}
}

func TestSuccessAfterTemp(t *testing.T) {
	tempErr := fmt.Errorf("this is a temp error")
	var count int
	fstore := failstore.NewFailstore(dstest.NewMapDatastoreForTest(t), func(op string) error {
		count++
		if count < 3 {
			return tempErr
		}
		count = 0
		return nil
	})

	rds := &Datastore{
		Batching: fstore,
		Retries:  5,
		TempErrFunc: func(err error) bool {
			return err == tempErr
		},
	}

	k := key.NewStrKey("test")
	val := []byte("foo")

	err := rds.Put(k, val)
	if err != nil {
		t.Fatal(err)
	}

	has, err := rds.Has(k)
	if err != nil {
		t.Fatal(err)
	}

	if !has {
		t.Fatal("should have this thing")
	}

	out, err := rds.Get(k)
	if err != nil {
		t.Fatal(err)
	}

	if string(out) != string(val) {
		t.Fatal("got wrong value")
	}
}
