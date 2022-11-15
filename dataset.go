/*
Copyright © 2022 James Darcy <jamesd@icr.ac.uk>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/

package dicom

import (
	"fmt"
	"sort"
	"strings"

	"github.com/JamesDarcy616/dicom/dcmerr"
	"github.com/JamesDarcy616/dicom/vr"
)

type Dataset struct {
	elems map[uint32]*Element
}

func NewDataset() *Dataset {
	ds := Dataset{elems: make(map[uint32]*Element)}
	return &ds
}

func (ds *Dataset) Get(tag uint32) (*Element, error) {
	elem, ok := ds.elems[tag]
	if !ok {
		return nil, dcmerr.Errorf(dcmerr.ErrNotFound, "Element 0x%08x not found", tag)
	}
	return elem, nil
}

func (ds *Dataset) GetString(tag uint32) (string, error) {
	elem, ok := ds.elems[tag]
	if !ok {
		return "", dcmerr.Errorf(dcmerr.ErrNotFound, "Element 0x%08x not found", tag)
	}
	value := elem.Value.Get()
	switch value := value.(type) {
	case string:
		return value, nil
	default:
		return "", dcmerr.Errorf(dcmerr.ErrNotFound, "Cannot convert element 0x%08x to string", tag)
	}
}

func (ds *Dataset) Iterator() DSIterator {
	return newIterator(ds)
}

func (ds *Dataset) Put(elem *Element) {
	if elem == nil {
		return
	}
	ds.elems[elem.Tag] = elem
}

func (ds *Dataset) PutString(tag uint32, vrStr, str string) error {
	if !vr.IsStringVR(vrStr) {
		return dcmerr.Errorf(dcmerr.ErrNotConvertible,
			"VR %v is not a string VR", vrStr)
	}
	value, err := NewValue(str)
	if err != nil {
		return err
	}
	elem := NewElement(tag, vrStr, uint32(len(str)), value)
	ds.elems[elem.Tag] = elem
	return nil
}

func (ds *Dataset) Size() int {
	return len(ds.elems)
}

func (ds *Dataset) String() string {
	str := ds.string("")
	if str[len(str)-1:] == "\n" {
		return str[:len(str)-1]
	}
	return str
}

func (ds *Dataset) string(indent string) string {
	var sb strings.Builder
	sb.Grow(ds.Size() * 128)
	iter := ds.Iterator()
	for iter.Next() {
		elem := iter.Value()
		sb.WriteString(fmt.Sprintf("%v%v\n", indent, elem))
		switch sq := elem.Value.(type) {
		case *sqValue:
			for _, item := range sq.value {
				sb.WriteString(item.string(">" + indent))
			}
		default:
		}
	}
	if err := iter.Err(); err != nil {
		fmt.Println(err.Error())
	}

	return sb.String()
}

type DSIterator interface {
	Next() bool
	Value() *Element
	Err() error
}

type iterator struct {
	ds      *Dataset
	keys    []uint32
	currIdx int
	err     error
}

func newIterator(ds *Dataset) DSIterator {
	keys := make([]uint32, ds.Size())
	idx := 0
	for key := range ds.elems {
		keys[idx] = key
		idx++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	// Create iterator with dataset, sorted key and index set before first key
	return &iterator{ds: ds, keys: keys, currIdx: -1}
}

func (i *iterator) Err() error {
	return i.err
}

func (i *iterator) Next() bool {
	if i.err != nil {
		return false
	}
	i.currIdx++
	return i.currIdx < len(i.keys)
}

func (i *iterator) Value() *Element {
	if i.currIdx >= len(i.keys) {
		i.err = dcmerr.NewErrIterInvalid()
		return nil
	}
	elem, ok := i.ds.elems[i.keys[i.currIdx]]
	if !ok {
		i.err = dcmerr.Errorf(dcmerr.ErrNotFound, "no element found for tag 0x%08x", i.keys[i.currIdx])
		return nil
	}
	return elem
}
