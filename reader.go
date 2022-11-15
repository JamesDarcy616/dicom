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

package dcm

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/JamesDarcy616/dicom/dcmerr"
)

type Reader interface {
	io.Reader
	ByteOrder() binary.ByteOrder
	BytesRead() uint64
	IsExplicit() bool
	Peek(n int) ([]byte, error)
	ReadFloat32() (float32, error)
	ReadFloat64() (float64, error)
	ReadInt16() (int16, error)
	ReadInt32() (int32, error)
	ReadString(n uint32) (string, error)
	ReadUint16() (uint16, error)
	ReadUint16LE() (uint16, error)
	ReadUint32() (uint32, error)
	ReadUint32LE() (uint32, error)
	SetExplicit(bool)
	Skip(n int64) error
}

type reader struct {
	explicit bool
	in       bufio.Reader
	nRead    uint64
	order    binary.ByteOrder
}

func NewReader(r io.Reader, order binary.ByteOrder, explicit bool) Reader {
	return &reader{
		explicit: explicit,
		in:       *bufio.NewReader(r),
		nRead:    0,
		order:    order,
	}
}

func (r *reader) ByteOrder() binary.ByteOrder {
	return r.order
}

func (r *reader) BytesRead() uint64 {
	return r.nRead
}

func (r *reader) IsExplicit() bool {
	return r.explicit
}

func (r *reader) Peek(n int) ([]byte, error) {
	return r.in.Peek(n)
}

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.in.Read(buf)
	// Increment nRead here as all other Read*() fns use this fn underneath
	r.nRead += uint64(n)
	if err != nil {
		return n, err
	}
	return n, err
}

func (r *reader) ReadFloat32() (float32, error) {
	var v float32
	if err := binary.Read(r, r.order, &v); err != nil {
		if err == io.EOF {
			return v, dcmerr.NewErrEOF()
		}
		if err == io.ErrUnexpectedEOF {
			return v, dcmerr.Errorf(dcmerr.ErrUnexpectedEOF,
				"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
		}
		return 0, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
	}
	return v, nil
}

func (r *reader) ReadFloat64() (float64, error) {
	var v float64
	if err := binary.Read(r, r.order, &v); err != nil {
		if err == io.EOF {
			return v, dcmerr.NewErrEOF()
		}
		if err == io.ErrUnexpectedEOF {
			return v, dcmerr.Errorf(dcmerr.ErrUnexpectedEOF,
				"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
		}
		return 0, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
	}
	return v, nil
}

func (r *reader) ReadInt16() (int16, error) {
	var v int16
	if err := binary.Read(r, r.order, &v); err != nil {
		if err == io.EOF {
			return v, dcmerr.NewErrEOF()
		}
		if err == io.ErrUnexpectedEOF {
			return v, dcmerr.Errorf(dcmerr.ErrUnexpectedEOF,
				"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
		}
		return 0, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
	}
	return v, nil
}

func (r *reader) ReadInt32() (int32, error) {
	var v int32
	if err := binary.Read(r, r.order, &v); err != nil {
		if err == io.EOF {
			return v, dcmerr.NewErrEOF()
		}
		if err == io.ErrUnexpectedEOF {
			return v, dcmerr.Errorf(dcmerr.ErrUnexpectedEOF,
				"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
		}
		return 0, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
	}
	return v, nil
}

func (r *reader) ReadString(len uint32) (string, error) {
	b := make([]byte, len)
	n, err := io.ReadFull(r, b)
	if uint32(n) != len {
		return "", dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - string not fully read", r.nRead, r.nRead)
	}
	if err != nil {
		return "", dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
	}
	return string(b), nil
}

func (r *reader) ReadUint16() (uint16, error) {
	var v uint16
	if err := binary.Read(r, r.order, &v); err != nil {
		if err == io.EOF {
			return v, dcmerr.NewErrEOF()
		}
		if err == io.ErrUnexpectedEOF {
			return v, dcmerr.Errorf(dcmerr.ErrUnexpectedEOF,
				"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
		}
		return 0, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
	}
	return v, nil
}

func (r *reader) ReadUint16LE() (uint16, error) {
	var v uint16
	if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
		if err == io.EOF {
			return v, dcmerr.NewErrEOF()
		}
		if err == io.ErrUnexpectedEOF {
			return v, dcmerr.Errorf(dcmerr.ErrUnexpectedEOF,
				"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
		}
		return 0, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
	}
	return v, nil
}

func (r *reader) ReadUint32() (uint32, error) {
	var v uint32
	if err := binary.Read(r, r.order, &v); err != nil {
		if err == io.EOF {
			return v, dcmerr.NewErrEOF()
		}
		if err == io.ErrUnexpectedEOF {
			return v, dcmerr.Errorf(dcmerr.ErrUnexpectedEOF,
				"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
		}
		return 0, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
	}
	return v, nil
}

func (r *reader) ReadUint32LE() (uint32, error) {
	var v uint32
	if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
		if err == io.EOF {
			return v, dcmerr.NewErrEOF()
		}
		if err == io.ErrUnexpectedEOF {
			return v, dcmerr.Errorf(dcmerr.ErrUnexpectedEOF,
				"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
		}
		return 0, dcmerr.Errorf(dcmerr.ErrIO,
			"error near byte %v (%08x) - %v", r.nRead, r.nRead, err.Error())
	}
	return v, nil
}

func (r *reader) SetExplicit(explicit bool) {
	r.explicit = explicit
}

func (r *reader) Skip(n int64) error {
	_, err := io.CopyN(io.Discard, r, n)
	if err != nil {
		dcmerr.Errorf(dcmerr.ErrIO,
			"error skipping ahead at byte %v (%08x) - %v", r.BytesRead(), r.BytesRead(), err.Error())
	}
	return nil
}
