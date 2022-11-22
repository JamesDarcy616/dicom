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

package dcmerr

import "fmt"

const (
	ErrNotFound int = iota
	ErrNotConvertible
	ErrTagNotFound
	ErrUIDNotFound

	ErrEOF
	ErrBadMagic
	ErrIO
	ErrSQItemFound
	ErrSQItemDelimFound
	ErrSQDelimFound
	ErrUnexpectedEOF

	ErrIterInvalid
	ErrUnsupported
	ErrNotImplemented
)

func IsErrNotFound(err error) bool {
	switch err := err.(type) {
	case DicomError:
		return err.Code() == ErrNotFound
	default:
		return false
	}
}

func IsErrNotConvertible(err error) bool {
	switch err := err.(type) {
	case DicomError:
		return err.Code() == ErrNotConvertible
	default:
		return false
	}
}

func IsErrSQItemDelimFound(err error) bool {
	switch err := err.(type) {
	case DicomError:
		return err.Code() == ErrSQItemDelimFound
	default:
		return false
	}
}

func IsErrTagNotFound(err error) bool {
	switch err := err.(type) {
	case DicomError:
		return err.Code() == ErrTagNotFound
	default:
		return false
	}
}

func IsErrUIDNotFound(err error) bool {
	switch err := err.(type) {
	case DicomError:
		return err.Code() == ErrUIDNotFound
	default:
		return false
	}
}

func IsErrEOF(err error) bool {
	switch err := err.(type) {
	case DicomError:
		return err.Code() == ErrEOF
	default:
		return false
	}
}

func IsErrBadMagic(err error) bool {
	switch err := err.(type) {
	case DicomError:
		return err.Code() == ErrBadMagic
	default:
		return false
	}
}

func IsErrIterInvalid(err error) bool {
	switch err := err.(type) {
	case DicomError:
		return err.Code() == ErrIterInvalid
	default:
		return false
	}
}

func NewErrEOF() DicomError {
	return &dicomError{msg: "EOF", code: ErrEOF}
}

func NewErrIterInvalid() DicomError {
	return &dicomError{msg: "IterInvalid", code: ErrIterInvalid}
}

func NewErrSQItemDelimFound() DicomError {
	return &dicomError{msg: "SQItemDelimFound", code: ErrSQItemDelimFound}
}

func NewErrUnexpectedEOF() DicomError {
	return &dicomError{msg: "UnexpectedEOF", code: ErrUnexpectedEOF}
}

type DicomError interface {
	error
	Code() int
}

type dicomError struct {
	msg  string
	code int
}

func NewDicomError(code int, msg string) DicomError {
	return &dicomError{msg: msg, code: code}
}

func (e *dicomError) Code() int {
	return e.code
}

func (e *dicomError) Error() string {
	return e.msg
}

func Errorf(code int, format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	return NewDicomError(code, msg)
}
