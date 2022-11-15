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
	"strings"

	"github.com/JamesDarcy616/dicom/tag"
)

type Element struct {
	Tag   uint32
	VR    string
	VL    uint32
	Value Value
}

func NewElement(tag uint32, vr string, vl uint32, value Value) *Element {
	return &Element{
		Tag:   tag,
		VR:    vr,
		VL:    vl,
		Value: value,
	}
}

func (e *Element) String() string {
	var sb strings.Builder
	sb.Grow(200)
	sb.WriteString(fmt.Sprintf("(%08x) %v #", e.Tag, e.VR))
	if e.VR == "SQ" {
		sq := e.Value.(*sqValue)
		n := len(sq.value)
		switch n {
		case 0:
			sb.WriteString("0 [] ")
		case 1:
			sb.WriteString("-1 [1 item] ")
		default:
			sb.WriteString(fmt.Sprintf("-1 [%v items] ", n))
		}
		sb.WriteString(fmt.Sprintf("%v", tag.Name(e.Tag)))
		return sb.String()
	}
	sb.WriteString(fmt.Sprintf("%v [%v] %v", e.VL, e.formatValue(64), tag.Name(e.Tag)))
	return sb.String()
}

func (e *Element) formatValue(max int) string {
	var sb strings.Builder
	sb.Grow(64)

	switch value := e.Value.(type) {
	case *stringValue:
		sb.WriteString(value.String())
	case *int16Value:
		last := len(value.value) - 1
		for i, v := range value.value {
			sb.WriteString(fmt.Sprint(v))
			if sb.Len() > 64 {
				break
			}
			if i < last {
				sb.WriteString("\\")
			}
		}
	case *uint16Value:
		last := len(value.value) - 1
		for i, v := range value.value {
			sb.WriteString(fmt.Sprint(v))
			if sb.Len() > 64 {
				break
			}
			if i < last {
				sb.WriteString("\\")
			}
		}
	case *int32Value:
		last := len(value.value) - 1
		for i, v := range value.value {
			sb.WriteString(fmt.Sprint(v))
			if sb.Len() > 64 {
				break
			}
			if i < last {
				sb.WriteString("\\")
			}
		}
	case *uint32Value:
		last := len(value.value) - 1
		for i, v := range value.value {
			sb.WriteString(fmt.Sprint(v))
			if sb.Len() > 64 {
				break
			}
			if i < last {
				sb.WriteString("\\")
			}
		}
	case *bytesValue:
		last := len(value.value) - 1
		for i, v := range value.value {
			sb.WriteString(fmt.Sprint(v))
			if sb.Len() > 64 {
				break
			}
			if i < last {
				sb.WriteString("\\")
			}
		}
	case *float32Value:
		last := len(value.value) - 1
		for i, v := range value.value {
			sb.WriteString(fmt.Sprint(v))
			if sb.Len() > 64 {
				break
			}
			if i < last {
				sb.WriteString("\\")
			}
		}
	case *float64Value:
		last := len(value.value) - 1
		for i, v := range value.value {
			sb.WriteString(fmt.Sprint(v))
			if sb.Len() > 64 {
				break
			}
			if i < last {
				sb.WriteString("\\")
			}
		}
	default:
		sb.WriteString(fmt.Sprintf("UNSUPPORTED VALUE TYPE %T", value))
	}

	str := sb.String()
	if len(str) > 64 {
		str = str[:60] + "..."
	}
	return str
}
