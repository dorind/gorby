package gorby

// Copyright(c) Dorin Duminica. All rights reserved.
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//   1. Redistributions of source code must retain the above copyright notice,
// 	 this list of conditions and the following disclaimer.
//
//   2. Redistributions in binary form must reproduce the above copyright notice,
// 	 this list of conditions and the following disclaimer in the documentation
// 	 and/or other materials provided with the distribution.
//
//   3. Neither the name of the copyright holder nor the names of its
// 	 contributors may be used to endorse or promote products derived from this
// 	 software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

import (
	"errors"
	"io"
	"unicode/utf8"
)

const (
	KRUNE_BUFF_MIN = 32 // 32 bytes is a nice sweet spot ~1/2 cpu cache line
)

var ErrBounds = errors.New("Rune index out of bounds!")
var ErrEncSize = errors.New("Invalid rune size in utf8.EncodeRune")
var ErrWriteSize = errors.New("io.writer implementor returned wrong number of written bytes")

type RuneBuff struct {
	// actual buffer
	buff []rune
	// index to next rune in buff
	// if it overflows, the buffer will automagically grow
	index uint
	// current buffer capacity
	capacity uint
}

// returns a new RuneBuffer instance with initial buffer capacity
// capacity will be set to KRUNE_BUFF_MIN if it's less than KRUNE_BUFF_MIN
func NewRuneBuff(capacity uint) *RuneBuff {
	if capacity < KRUNE_BUFF_MIN {
		capacity = KRUNE_BUFF_MIN
	}
	r := &RuneBuff{
		index:    0,
		capacity: capacity,
	}
	r.Reset()
	return r
}

// reset state
func (m *RuneBuff) Reset() {
	clen := uint(len(m.buff))
	if clen < m.capacity {
		m.buff = make([]rune, m.capacity)
	}
	m.index = 0
}

// returns true if the buffer contains (r) rune
func (m *RuneBuff) Contains(r rune) bool {
	return m.IndexOf(r) >= 0
}

// we know for sure that buff has at least KRUNE_BUFF_MIN runes allocated
func (m *RuneBuff) First() rune {
	return m.buff[0]
}

// returns index of the (r) rune in buffer
// if not found, will return (-1)
func (m *RuneBuff) IndexOf(r rune) int {
	for i := uint(0); i < m.index; i++ {
		if m.buff[i] == r {
			return int(i)
		}
	}
	return -1
}

// rune at (index) in buffer
// errors
//	ErrRuneIndexOutOfBounds
func (m *RuneBuff) AtIndex(index uint) (rune, error) {
	if index >= m.index {
		return 0, ErrBounds
	}

	return m.buff[index], nil
}

// returns the current buffer as a string
// if there's no rune set in buffer, it returns an empty string
func (m *RuneBuff) String() string {
	s := ""
	i := m.index
	if i > 0 {
		s = string(m.buff[:i])
	}
	return s
}

// grow the internal buffer to (capacity)
func (m *RuneBuff) grow(capacity uint) {
	m.capacity = capacity
	t := make([]rune, m.capacity)
	copy(t, m.buff)
	m.buff = t
}

func (m *RuneBuff) autoGrow() {
	// increase buffer capacity by 1/4th
	ncap := m.capacity + (m.capacity >> 2)
	m.grow(ncap)
}

// add (r) rune to buffer
// if needed, buffer will automatically grow
func (m *RuneBuff) PushRune(r rune) {
	if m.index >= m.capacity {
		m.autoGrow()
	}
	m.buff[m.index] = r
	m.index++
}

// append string to buffer
func (m *RuneBuff) PushString(s string) {
	// extract runes from string
	srunes := []rune(s)
	// length never changes, store it
	l := len(srunes)
	// loop over runes and copy to buffer
	for i := 0; i < l; i++ {
		if m.index >= m.capacity {
			m.autoGrow()
		}
		m.buff[m.index] = srunes[i]
		m.index++
	}
}

// returns current index
func (m *RuneBuff) Index() uint {
	return m.index
}

// sets index to new value if it doesn't overflow
func (m *RuneBuff) SetIndex(index uint) {
	if index < m.capacity {
		m.index = index
	}
}

// returns the capacity of the buffer
func (m *RuneBuff) Capacity() uint {
	return m.capacity
}

// write all runes in buffer to writer as utf8 encoded
func (m *RuneBuff) WriteTo(writer io.Writer) (int, error) {
	// grab buff ref locally
	buff := m.buff
	// allocate buffer for rune encoding
	wbuf := make([]byte, utf8.UTFMax)
	// total bytes
	tbw := 0
	// loop over all runes in buff
	for i := uint(0); i < m.index; i++ {
		// fetch rune
		r := buff[i]
		// encode rune into wbuf
		size := utf8.EncodeRune(wbuf, r)
		if (size < 1) || (size > utf8.UTFMax) {
			// should never happen
			return tbw, ErrEncSize
		}
		// write rune to wrier
		bw, err := writer.Write(wbuf[:size])
		if err != nil {
			// rats! error while trying to write
			return tbw, err
		}
		if bw != size {
			// this should *technically* never happen
			return tbw, ErrWriteSize
		}
		// keep track of written bytes
		tbw += bw
	}
	// everything went great, return total byte count
	return tbw, nil
}
