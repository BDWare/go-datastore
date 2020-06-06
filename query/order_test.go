// Copyright for portions of this fork are held by [Juan Batiz-Benet, 2016] as
// part of the original go-datastore project. All other copyright for
// this fork are held by [The BDWare Authors, 2020]. All rights reserved.
// Use of this source code is governed by MIT license that can be
// found in the LICENSE file.

package query

import (
	"testing"

	key "github.com/bdware/go-datastore/key"
)

func testKeyOrder(t *testing.T, f Order, keys []key.Key, expect key.KeySlice) {
	t.Helper()

	e := make([]Entry, len(keys))
	for i, k := range keys {
		e[i] = Entry{Key: k}
	}

	res := ResultsWithEntries(Query{}, e)
	res = NaiveOrder(res, f)
	actualE, err := res.Rest()
	if err != nil {
		t.Fatal(err)
	}

	actual := key.KeySlice(make([]key.Key, len(actualE)))
	for i, e := range actualE {
		actual[i] = e.Key
	}

	if len(actual) != len(expect) {
		t.Error("expect != actual.", expect, actual)
	}

	if !actual.Join().Equal(expect.Join()) {
		t.Error("expect != actual.", expect, actual)
	}
}

func TestOrderByKey(t *testing.T) {

	testKeyOrder(t, OrderByKey{}, sampleKeys, key.StrsToKeys([]string{
		"/a",
		"/ab",
		"/ab/c",
		"/ab/cd",
		"/ab/ef",
		"/ab/fg",
		"/abce",
		"/abcf",
	}))
	testKeyOrder(t, OrderByKeyDescending{}, sampleKeys, key.StrsToKeys([]string{
		"/abcf",
		"/abce",
		"/ab/fg",
		"/ab/ef",
		"/ab/cd",
		"/ab/c",
		"/ab",
		"/a",
	}))
}
