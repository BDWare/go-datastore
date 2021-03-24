// Copyright for portions of this fork are held by [Juan Batiz-Benet, 2016]
// as part of the original go-datastore project. All other copyright for this
// fork are held by [DAOT Labs, 2020]. All rights reserved. Use of this source
// code is governed by MIT license that can be found in the LICENSE file.

// Package key provides the Key interface and the KeySlice type, along with
// some utility functions around key.
package key

import (
	"errors"
	"fmt"
)

type KeyType uint8

const (
	// Key backed by string
	KeyTypeString KeyType = iota
	// Key backed by byte slice
	KeyTypeBytes
)

var (
	ErrKeyTypeNotSupported = errors.New("key type not supported")
)

/*
A Key represents the unique identifier of an object.
Keys are meant to be unique across a system.
There are two Key implementations:
StrKey backed by string and BytesKey backed by byte slice.
*/
type Key interface {
	fmt.Stringer
	// KeyType returns the key type
	KeyType() KeyType
	// Bytes returns the string value of Key as a []byte
	Bytes() []byte
	// Equal checks equality of two keys
	Equal(k2 Key) bool
	// Less checks whether this key is sorted lower than another.
	Less(k2 Key) bool
	// Child returns the `child` Key of this Key.
	Child(k2 Key) Key
	// IsAncestorOf returns whether this key is a prefix of `other`
	IsAncestorOf(other Key) bool
	// IsDescendantOf returns whether this key contains another as a prefix (excluding equals).
	IsDescendantOf(other Key) bool
	// HasPrefix returns whether this key contains another as a prefix (including equals).
	HasPrefix(prefix Key) bool
	// HasPrefix returns whether this key contains another as a suffix (including equals).
	HasSuffix(suffix Key) bool
	// TrimPrefix returns a new key equals to this key without the provided leading prefix key.
	// If s doesn't start with prefix, this key is returned unchanged.
	TrimPrefix(prefix Key) Key
	// TrimSuffix returns a new key equals to this key without the provided trailing suffix key.
	// If s doesn't end with suffix, this key is returned unchanged.
	TrimSuffix(suffix Key) Key
	// MarshalJSON implements the json.Marshaler interface,
	// keys are represented as JSON strings
	MarshalJSON() ([]byte, error)
}

// Clean up a StrKey, using path.Clean, no-op for BytesKey.
func Clean(k Key) Key {
	if k == nil {
		return nil
	}
	switch k.KeyType() {
	case KeyTypeString:
		sk := k.(StrKey)
		sk.Clean()
		return sk
	case KeyTypeBytes:
		return k
	default:
		panic(ErrKeyTypeNotSupported)
	}
}

// Compare returns an integer comparing two Keys lexicographically.
// The result will be 0 if a.Equal(b), -1 if a.Less(b), and +1 if b.Less(a).
func Compare(a, b Key) int {
	if a == nil {
		if b == nil {
			return 0
		} else {
			return -1
		}
	}
	if a.Equal(b) {
		return 0
	}
	if a.Less(b) {
		return -1
	}
	return +1
}

// KeySlice attaches the methods of sort.Interface to []Key,
// sorting in increasing order.
type KeySlice []Key

func (p KeySlice) Len() int           { return len(p) }
func (p KeySlice) Less(i, j int) bool { return p[i].Less(p[j]) }
func (p KeySlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Join Joins keys in the KeySlice into a single key,
// returns NewStrKey("") if slice is empty.
func (p KeySlice) Join() Key {
	if len(p) == 0 {
		return NewStrKey("")
	}
	key := p[0]
	for _, k := range p[1:] {
		key = key.Child(k)
	}
	return key
}
