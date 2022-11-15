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

package vr

const (
	AE string = "AE"
	AS string = "AS"
	AT string = "AT"
	CS string = "CS"
	DA string = "DA"
	DS string = "DS"
	DT string = "DT"
	FD string = "FD"
	FL string = "FL"
	IS string = "IS"
	LO string = "LO"
	LT string = "LT"
	OB string = "OB"
	OD string = "OD"
	OF string = "OF"
	OL string = "OL"
	OV string = "OV"
	OW string = "OW"
	PN string = "PN"
	SH string = "SH"
	SL string = "SL"
	SQ string = "SQ"
	SS string = "SS"
	ST string = "ST"
	SV string = "SV"
	TM string = "TM"
	UC string = "UC"
	UI string = "UI"
	UL string = "UL"
	UN string = "UN"
	UR string = "UR"
	US string = "US"
	UT string = "UT"
	UV string = "UV"
)

var strVRs = make(map[string]struct{})

func IsStringVR(value string) bool {
	_, ok := strVRs[value]
	return ok
}

func init() {
	if len(strVRs) == 0 {
		initStrVRs()
	}
}

func initStrVRs() {
	strVRs[AE] = struct{}{}
	strVRs[AS] = struct{}{}
	strVRs[CS] = struct{}{}
	strVRs[DA] = struct{}{}
	strVRs[DS] = struct{}{}
	strVRs[DT] = struct{}{}
	strVRs[IS] = struct{}{}
	strVRs[LO] = struct{}{}
	strVRs[LT] = struct{}{}
	strVRs[PN] = struct{}{}
	strVRs[SH] = struct{}{}
	strVRs[ST] = struct{}{}
	strVRs[TM] = struct{}{}
	strVRs[UC] = struct{}{}
	strVRs[UI] = struct{}{}
	strVRs[UR] = struct{}{}
	strVRs[UT] = struct{}{}
}
