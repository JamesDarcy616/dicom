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

package tag

var tagMap = make(map[uint32]*TagInfo)

func Name(tag uint32) string {
	info, ok := tagMap[tag]
	if !ok {
		return "UNKNOWN"
	}
	return info.Name()
}

func VR(tag uint32) string {
	info, ok := tagMap[tag]
	if !ok {
		return "UN"
	}
	return info.VR()
}

type TagInfo struct {
	tag  uint32
	vr   string
	name string
	vm   string
}

func NewTagInfo(tag uint32, vr, name, vm string) *TagInfo {
	return &TagInfo{tag: tag, vr: vr, name: name, vm: vm}
}

func (ti *TagInfo) Name() string {
	return ti.name
}

func (ti *TagInfo) Tag() uint32 {
	return ti.tag
}

func (ti *TagInfo) VM() string {
	return ti.vm
}

func (ti *TagInfo) VR() string {
	return ti.vr
}
