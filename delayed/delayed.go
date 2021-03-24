// Copyright for portions of this fork are held by [Juan Batiz-Benet, 2016]
// as part of the original go-datastore project. All other copyright for this
// fork are held by [DAOT Labs, 2020]. All rights reserved. Use of this source
// code is governed by MIT license that can be found in the LICENSE file.

// Package delayed wraps a datastore allowing to artificially
// delay all operations.
package delayed

import (
	"io"

	ds "github.com/daotl/go-datastore"
	key "github.com/daotl/go-datastore/key"
	dsq "github.com/daotl/go-datastore/query"
	delay "github.com/ipfs/go-ipfs-delay"
)

// New returns a new delayed datastore.
func New(ds ds.Datastore, delay delay.D) *Delayed {
	return &Delayed{ds: ds, delay: delay}
}

// Delayed is an adapter that delays operations on the inner datastore.
type Delayed struct {
	ds    ds.Datastore
	delay delay.D
}

var _ ds.Batching = (*Delayed)(nil)
var _ io.Closer = (*Delayed)(nil)

// Put implements the ds.Datastore interface.
func (dds *Delayed) Put(key key.Key, value []byte) (err error) {
	dds.delay.Wait()
	return dds.ds.Put(key, value)
}

// Sync implements Datastore.Sync
func (dds *Delayed) Sync(prefix key.Key) error {
	dds.delay.Wait()
	return dds.ds.Sync(prefix)
}

// Get implements the ds.Datastore interface.
func (dds *Delayed) Get(key key.Key) (value []byte, err error) {
	dds.delay.Wait()
	return dds.ds.Get(key)
}

// Has implements the ds.Datastore interface.
func (dds *Delayed) Has(key key.Key) (exists bool, err error) {
	dds.delay.Wait()
	return dds.ds.Has(key)
}

// GetSize implements the ds.Datastore interface.
func (dds *Delayed) GetSize(key key.Key) (size int, err error) {
	dds.delay.Wait()
	return dds.ds.GetSize(key)
}

// Delete implements the ds.Datastore interface.
func (dds *Delayed) Delete(key key.Key) (err error) {
	dds.delay.Wait()
	return dds.ds.Delete(key)
}

// Query implements the ds.Datastore interface.
func (dds *Delayed) Query(q dsq.Query) (dsq.Results, error) {
	dds.delay.Wait()
	return dds.ds.Query(q)
}

// Batch implements the ds.Batching interface.
func (dds *Delayed) Batch() (ds.Batch, error) {
	return ds.NewBasicBatch(dds), nil
}

// DiskUsage implements the ds.PersistentDatastore interface.
func (dds *Delayed) DiskUsage() (uint64, error) {
	dds.delay.Wait()
	return ds.DiskUsage(dds.ds)
}

// Close closes the inner datastore (if it implements the io.Closer interface).
func (dds *Delayed) Close() error {
	if closer, ok := dds.ds.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
