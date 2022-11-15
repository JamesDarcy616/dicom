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
	"encoding/binary"
	"io"
	"os"
	"sync"

	"github.com/JamesDarcy616/dicom/dcmerr"
	"github.com/JamesDarcy616/dicom/tag"
	"github.com/JamesDarcy616/dicom/uid"
)

const magic = "DICM"
const (
	UndefinedLength = 0xffffffff
	SQItem          = 0xfffee000
	SQItemDelim     = 0xfffee00d
	SQDelim         = 0xfffee0dd
)

type Parser struct {
	// Protect internal state during a call
	mutex sync.Mutex
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(r io.Reader, explicit bool) (*Dataset, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	ds := NewDataset()
	reader := NewReader(r, binary.LittleEndian, explicit)
	err := p.parseAll(ds, reader)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (p *Parser) ParseFile(filename string) (*Dataset, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Start in ExpLE mode for file metadata
	reader := NewReader(file, binary.LittleEndian, true)

	if err := p.checkHeader(reader); err != nil {
		return nil, err
	}

	ds := NewDataset()
	err = p.parseFileMeta(ds, reader)
	if err != nil {
		return nil, err
	}
	if err := p.checkXferSyntax(ds, reader); err != nil {
		return nil, err
	}
	err = p.parseAll(ds, reader)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (p *Parser) ParseFileUntil(filename string, maxTag uint32) (*Dataset, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Start in ExpLE mode for file metadata
	reader := NewReader(file, binary.LittleEndian, true)

	if err := p.checkHeader(reader); err != nil {
		return nil, err
	}

	ds := NewDataset()
	err = p.parseFileMeta(ds, reader)
	if err != nil {
		return nil, err
	}
	if err := p.checkXferSyntax(ds, reader); err != nil {
		return nil, err
	}
	err = p.parseUntil(ds, reader, maxTag)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (p *Parser) ParseUntil(r io.Reader, explicit bool, maxTag uint32) (*Dataset, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	ds := NewDataset()
	reader := NewReader(r, binary.LittleEndian, explicit)
	err := p.parseUntil(ds, reader, maxTag)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (p *Parser) checkHeader(r Reader) error {
	preamble := make([]byte, 128)
	n, err := r.Read(preamble)
	if n != len(preamble) {
		return dcmerr.Errorf(dcmerr.ErrIO,
			"error reading preamble at byte %v (%08x)", r.BytesRead(), r.BytesRead())
	}
	if err != nil {
		return err
	}
	prefix := make([]byte, 4)
	n, err = r.Read(prefix)
	if n != len(prefix) || err != nil {
		return err
	}
	if string(prefix) != magic {
		return dcmerr.Errorf(dcmerr.ErrIO,
			"bad magic %v at byte %v (%08x)", prefix, r.BytesRead(), r.BytesRead())
	}
	return nil
}

func (p *Parser) checkSQMarker(tag uint32) (*Element, error) {
	value, err := NewValue(nil)
	if err != nil {
		return nil, err
	}
	switch tag {
	case SQItem:
		return NewElement(tag, "", UndefinedLength, value), nil
	case SQItemDelim:
		return NewElement(tag, "", UndefinedLength, value), nil
	case SQDelim:
		return NewElement(tag, "", UndefinedLength, value), nil
	}
	return nil, nil
}

func (p *Parser) checkXferSyntax(ds *Dataset, r Reader) error {
	tsuid, err := ds.GetString(tag.TransferSyntaxUID)
	if err != nil {
		return err
	}
	if tsuid == uid.ImplicitVRLittleEndian {
		r.SetExplicit(false)
	}
	return nil
}

func (p *Parser) parseAll(ds *Dataset, r Reader) error {
	for {
		elem, err := p.readElement(r)
		if err != nil {
			if dcmerr.IsErrEOF(err) {
				return nil
			}
			return err
		}
		// Detect end of dataset when inside an SQ
		if elem.Tag == SQItemDelim {
			return nil
		}
		ds.Put(elem)
	}
}

func (p *Parser) parseFileMeta(ds *Dataset, r Reader) error {
	start := r.BytesRead()
	metaLen, err := p.readElement(r)
	if err != nil {
		return err
	}
	ds.Put(metaLen)
	maxRead := start + uint64(metaLen.Value.Get().(uint32))
	for r.BytesRead() < maxRead {
		elem, err := p.readElement(r)
		if err != nil {
			return err
		}
		ds.Put(elem)
	}
	return nil
}

func (p *Parser) parseUntil(ds *Dataset, r Reader, maxTag uint32) error {
	for {
		elem, err := p.readElementPeek(r, maxTag)
		if err != nil {
			if dcmerr.IsErrEOF(err) {
				return nil
			}
			return err
		}
		// Detect end of dataset when inside an SQ
		if elem.Tag == SQItemDelim {
			return nil
		}
		ds.Put(elem)
		// Bail out on reaching last required tag
		if elem.Tag == maxTag {
			return nil
		}
	}
}

func (p *Parser) readBytesValue(r Reader, vl uint32) (Value, error) {
	buf := make([]byte, vl)
	n, err := io.ReadFull(r, buf)
	if uint32(n) < vl {
		return nil, err
	}
	return NewValue(buf)
}

func (p *Parser) readElement(r Reader) (*Element, error) {
	tag, err := p.readLETag(r)
	if err != nil {
		return nil, err
	}
	// Short circuit on SQ markers
	sqMarker, err := p.checkSQMarker(tag)
	if err != nil {
		return nil, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
	}
	if sqMarker != nil {
		// Consume the remaining bytes of the SQ demarcation element
		if err := r.Skip(4); err != nil {
			return nil, err
		}
		return sqMarker, err
	}

	vr, err := p.readVR(r, tag)
	if err != nil {
		return nil, err
	}
	vl, err := p.readVL(r, vr)
	if err != nil {
		return nil, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
	}
	value, err := p.readValue(r, vr, vl)
	if err != nil {
		return nil, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
	}

	return NewElement(tag, vr, vl, value), nil
}

func (p *Parser) readElementPeek(r Reader, maxTag uint32) (*Element, error) {
	peek, err := r.Peek(4)
	if err != nil {
		if err == io.EOF {
			return nil, dcmerr.NewErrEOF()
		}
		return nil, dcmerr.Errorf(dcmerr.ErrIO,
			"error peeking at byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
	}
	tag := p.readLETagBytes(peek[0:4])
	if tag > maxTag {
		// Send ErrEOF to simulate the end of the stream
		return nil, dcmerr.NewErrEOF()
	}
	// Short circuit on SQ markers
	sqMarker, err := p.checkSQMarker(tag)
	if err != nil {
		return nil, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
	}
	if sqMarker != nil {
		// Consume the 8 bytes of the SQ demarcation element
		if err := r.Skip(8); err != nil {
			return nil,
				dcmerr.Errorf(dcmerr.ErrIO,
					"error skipping ahead at byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
		}
		return sqMarker, err
	}

	// Skip the 4 peeked bytes
	if err := r.Skip(4); err != nil {
		return nil,
			dcmerr.Errorf(dcmerr.ErrIO,
				"error skipping ahead at byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
	}

	vr, err := p.readVR(r, tag)
	if err != nil {
		return nil, err
	}
	vl, err := p.readVL(r, vr)
	if err != nil {
		return nil, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
	}
	value, err := p.readValue(r, vr, vl)
	if err != nil {
		return nil, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
	}

	return NewElement(tag, vr, vl, value), nil
}

// Read an element in ImpLE format regardless of TransferSyntax e.g. SQ items
func (p *Parser) readImpLEElement(r Reader) (*Element, error) {
	tag32, err := p.readLETag(r)
	if err != nil {
		return nil, err
	}
	vl, err := r.ReadUint32LE()
	if err != nil {
		return nil, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
	}
	// Short circuit on SQ markers
	sqMarker, err := p.checkSQMarker(tag32)
	if err != nil {
		return nil, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
	}
	if sqMarker != nil {
		return sqMarker, err
	}

	vr := tag.VR(tag32)
	value, err := p.readValue(r, vr, vl)
	if err != nil {
		return nil, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
	}

	return NewElement(tag32, vr, vl, value), nil
}

func (p *Parser) readFloat32Value(r Reader, vl uint32) (Value, error) {
	data := make([]float32, vl/4)
	for i := range data {
		v, err := r.ReadFloat32()
		if err != nil {
			return nil, err
		}
		data[i] = v
	}
	return NewValue(data)
}

func (p *Parser) readFloat64Value(r Reader, vl uint32) (Value, error) {
	data := make([]float64, vl/8)
	for i := range data {
		v, err := r.ReadFloat64()
		if err != nil {
			return nil, err
		}
		data[i] = v
	}
	return NewValue(data)
}

func (p *Parser) readInt16Value(r Reader, vl uint32) (Value, error) {
	data := make([]int16, vl/2)
	for i := range data {
		v, err := r.ReadInt16()
		if err != nil {
			return nil, err
		}
		data[i] = v
	}
	return NewValue(data)
}

func (p *Parser) readInt32Value(r Reader, vl uint32) (Value, error) {
	data := make([]int32, vl/4)
	for i := range data {
		v, err := r.ReadInt32()
		if err != nil {
			return nil, err
		}
		data[i] = v
	}
	return NewValue(data)
}

// Read a tag in LE ordering regardless of the TransferSyntax e.g. SQ items
func (p *Parser) readLETag(r Reader) (uint32, error) {
	g, err := r.ReadUint16LE()
	if err != nil {
		return 0, err
	}
	e, err := r.ReadUint16LE()
	if err != nil {
		// EOF should not happen except at the beginning of a tag
		if dcmerr.IsErrEOF(err) {
			return 0, dcmerr.Errorf(dcmerr.ErrUnexpectedEOF,
				"error near byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
		}
		return 0, err
	}
	return uint32(g)<<16 | uint32(e), nil
}

func (p *Parser) readLETagBytes(b []byte) uint32 {
	return uint32(b[0])<<16 | uint32(b[1])<<24 | uint32(b[2]) | uint32(b[3])<<8
}

func (p *Parser) readSequence(r Reader) (Value, error) {
	items := make([]*Dataset, 0)
	for {
		pos := r.BytesRead()
		elem, err := p.readImpLEElement(r)
		if err != nil {
			return nil, err
		}
		// Bail if the end of the SQ
		if elem.Tag == SQDelim {
			break
		}
		if elem.Tag != SQItem {
			return nil, dcmerr.Errorf(dcmerr.ErrIO,
				"SQItem tag expected at %v (%08x), found %08x", pos, pos, elem.Tag)
		}
		ds := NewDataset()
		if err := p.parseAll(ds, r); err != nil {
			return nil, err
		}
		items = append(items, ds)
	}
	return NewValue(items)
}

func (p *Parser) readStringValue(r Reader, vl uint32) (Value, error) {
	s, err := r.ReadString(vl)
	if err != nil {
		return nil, err
	}
	return NewValue(s)
}

func (p *Parser) readUint16Value(r Reader, vl uint32) (Value, error) {
	data := make([]uint16, vl/2)
	for i := range data {
		v, err := r.ReadUint16()
		if err != nil {
			return nil, err
		}
		data[i] = v
	}
	return NewValue(data)
}

func (p *Parser) readUint32Value(r Reader, vl uint32) (Value, error) {
	data := make([]uint32, vl/4)
	for i := range data {
		v, err := r.ReadUint32()
		if err != nil {
			return nil, err
		}
		data[i] = v
	}
	return NewValue(data)
}

func (p *Parser) readUndefLenValue(r Reader, vr string) (Value, error) {
	switch vr {
	case "SQ":
		// fmt.Printf("SQ found at %v (%08x)\n", r.BytesRead()-8, r.BytesRead()-8)
		return p.readSequence(r)
	}
	return nil, dcmerr.Errorf(dcmerr.ErrIO, "undefined length for VR %v", vr)
}

func (p *Parser) readValue(r Reader, vr string, vl uint32) (Value, error) {
	if vl == UndefinedLength {
		return p.readUndefLenValue(r, vr)
	}
	switch vr {
	case "AE", "AS", "CS", "DA", "DS", "DT", "IS", "LO", "LT", "PN", "SH", "ST",
		"TM", "UI":
		return p.readStringValue(r, vl)
	case "UL":
		return p.readUint32Value(r, vl)
	// "xs" (from dcmtk.dic) means "either US or SS", read as US - it can be converted if required later
	case "US", "xs":
		return p.readUint16Value(r, vl)
	// "ox", "px" (from dcmtk.dic) mean "either OB or OW", read as OB - it can be converted if required later
	case "OB", "UN", "ox", "px":
		return p.readBytesValue(r, vl)
	// case "AT":
	// 	use16 = true
	case "FL":
		return p.readFloat32Value(r, vl)
	case "FD":
		return p.readFloat64Value(r, vl)
	case "SL":
		return p.readInt32Value(r, vl)
	case "SS", "OW":
		return p.readInt16Value(r, vl)
	default:
		return nil, dcmerr.Errorf(dcmerr.ErrIO, "unsupported VR: %v", vr)
	}
}

func (p *Parser) readVL(r Reader, vr string) (uint32, error) {
	if !r.IsExplicit() {
		return r.ReadUint32()
	}
	var use16 bool
	switch vr {
	case "AE", "AS", "AT", "CS", "DA", "DS", "DT", "FL", "FD", "IS", "LO",
		"LT", "PN", "SH", "SL", "SS", "ST", "TM", "UI", "UL", "US":
		use16 = true
	}
	if use16 {
		vl16, err := r.ReadUint16()
		if err != nil {
			return 0, err
		}
		return uint32(vl16), nil
	}
	if _, err := r.ReadUint16(); err != nil {
		return 0, err
	}
	vl, err := r.ReadUint32()
	if err != nil {
		return 0, err
	}
	return vl, nil
}

func (p *Parser) readVR(r Reader, tag32 uint32) (string, error) {
	if !r.IsExplicit() {
		return tag.VR(tag32), nil
	}
	return r.ReadString(2)
}
