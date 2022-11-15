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

package uid

import (
	"xnatrest/internal/dcm/dcmerr"
)

var uidMap = make(map[string]*UIDInfo)

func Name(value string) (string, error) {
	info, ok := uidMap[value]
	if !ok {
		return "", dcmerr.Errorf(dcmerr.ErrUIDNotFound, "Unknown UID: %v", value)
	}
	return info.Name(), nil
}

type UIDInfo struct {
	name    string
	value   string
	desc    string
	uidType string
}

func NewUIDInfo(value, name, desc, uidType string) *UIDInfo {
	return &UIDInfo{value: value, name: name, desc: desc, uidType: uidType}
}

func (uid *UIDInfo) Description() string {
	return uid.desc
}

func (uid *UIDInfo) Name() string {
	return uid.name
}

func (uid *UIDInfo) Type() string {
	return uid.uidType
}

func (uid *UIDInfo) Value() string {
	return uid.value
}
