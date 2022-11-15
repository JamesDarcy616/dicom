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
	"fmt"
	"strings"

	"xnatrest/internal/dcm/dcmerr"
)

type Value interface {
	Get() interface{}
	GetAll() interface{}
	String() string
}

func NewValue(raw interface{}) (Value, error) {
	if raw == nil {
		return &emptyValue{}, nil
	}
	switch raw := raw.(type) {
	case string:
		return &stringValue{value: nullStrip(raw)}, nil
	case []uint32:
		return &uint32Value{value: raw}, nil
	case []uint16:
		return &uint16Value{value: raw}, nil
	case []int32:
		return &int32Value{value: raw}, nil
	case []int16:
		return &int16Value{value: raw}, nil
	case []byte:
		return &bytesValue{value: raw}, nil
	case []float32:
		return &float32Value{value: raw}, nil
	case []float64:
		return &float64Value{value: raw}, nil
	case []*Dataset:
		return &sqValue{value: raw}, nil
	default:
		return nil, dcmerr.Errorf(dcmerr.ErrUnsupported, "unknown or unsupported type: %T", raw)
	}
}

type bytesValue struct {
	value []byte
}

func (v *bytesValue) Get() interface{}    { return v.value }
func (v *bytesValue) GetAll() interface{} { return v.value }
func (v *bytesValue) String() string      { return fmt.Sprintf("%v", v.value) }

type emptyValue struct{}

func (v *emptyValue) Get() interface{}    { return nil }
func (v *emptyValue) GetAll() interface{} { return nil }
func (v *emptyValue) String() string      { return fmt.Sprintf("%v", nil) }

type float32Value struct {
	value []float32
}

func (v *float32Value) Get() interface{}    { return v.value[0] }
func (v *float32Value) GetAll() interface{} { return v.value }
func (v *float32Value) String() string      { return fmt.Sprintf("%v", v.value) }

type float64Value struct {
	value []float64
}

func (v *float64Value) Get() interface{}    { return v.value[0] }
func (v *float64Value) GetAll() interface{} { return v.value }
func (v *float64Value) String() string      { return fmt.Sprintf("%v", v.value) }

type int16Value struct {
	value []int16
}

func (v *int16Value) Get() interface{}    { return v.value[0] }
func (v *int16Value) GetAll() interface{} { return v.value }
func (v *int16Value) String() string      { return fmt.Sprintf("%v", v.value) }

type int32Value struct {
	value []int32
}

func (v *int32Value) Get() interface{}    { return v.value[0] }
func (v *int32Value) GetAll() interface{} { return v.value }
func (v *int32Value) String() string      { return fmt.Sprintf("%v", v.value) }

type stringValue struct {
	value string
}

func (v *stringValue) Get() interface{}    { return strings.TrimSpace(v.value) }
func (v *stringValue) GetAll() interface{} { return strings.TrimSpace(v.value) }
func (v *stringValue) String() string      { return strings.TrimSpace(v.value) }

type sqValue struct {
	value []*Dataset
}

func (v *sqValue) Get() interface{}    { return v.value[0] }
func (v *sqValue) GetAll() interface{} { return v.value }
func (v *sqValue) String() string      { return "" }

type uint16Value struct {
	value []uint16
}

func (v *uint16Value) Get() interface{}    { return v.value[0] }
func (v *uint16Value) GetAll() interface{} { return v.value }
func (v *uint16Value) String() string      { return fmt.Sprintf("%v", v.value) }

type uint32Value struct {
	value []uint32
}

func (v *uint32Value) Get() interface{}    { return v.value[0] }
func (v *uint32Value) GetAll() interface{} { return v.value }
func (v *uint32Value) String() string      { return fmt.Sprintf("%v", v.value) }

func nullStrip(in string) string {
	if len(in) < 2 {
		return in
	}
	if in[len(in)-1:] != "\x00" {
		return in
	}
	return in[:len(in)-1]
}
