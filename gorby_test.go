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
	"bytes"
	"testing"
)

var test_string1 = "abracadabra"
var test_string_09 = "0123456789"
var test_string_u1 = "単体テストは楽しいですが、ご覧のとおり、問題が発生する可能性のある場所はたくさんあります。問題が発生する可能性がある場合は、"

// the following samples are from
// http://www.columbia.edu/~fdc/utf8/index.html
// rune poem
var test_string_u2 = `ᚠᛇᚻ᛫ᛒᛦᚦ᛫ᚠᚱᚩᚠᚢᚱ᛫ᚠᛁᚱᚪ᛫ᚷᛖᚻᚹᛦᛚᚳᚢᛗ
ᛋᚳᛖᚪᛚ᛫ᚦᛖᚪᚻ᛫ᛗᚪᚾᚾᚪ᛫ᚷᛖᚻᚹᛦᛚᚳ᛫ᛗᛁᚳᛚᚢᚾ᛫ᚻᛦᛏ᛫ᛞᚫᛚᚪᚾ
ᚷᛁᚠ᛫ᚻᛖ᛫ᚹᛁᛚᛖ᛫ᚠᚩᚱ᛫ᛞᚱᛁᚻᛏᚾᛖ᛫ᛞᚩᛗᛖᛋ᛫ᚻᛚᛇᛏᚪᚾ᛬`

// Laȝamon'
var test_string_u3 = `An preost wes on leoden, Laȝamon was ihoten
He wes Leovenaðes sone -- liðe him be Drihten.
He wonede at Ernleȝe at æðelen are chirechen,
Uppen Sevarne staþe, sel þar him þuhte,
Onfest Radestone, þer he bock radde.`

// Tagelied
var test_string_u4 = `Sîne klâwen durh die wolken sint geslagen,
er stîget ûf mit grôzer kraft,
ich sih in grâwen tägelîch als er wil tagen,
den tac, der im geselleschaft
erwenden wil, dem werden man,
den ich mit sorgen în verliez.
ich bringe in hinnen, ob ich kan.
sîn vil manegiu tugent michz leisten hiez.`

// Odysseus Elytis - Monotonic
var test_string_u5 = `Τη γλώσσα μου έδωσαν ελληνική
το σπίτι φτωχικό στις αμμουδιές του Ομήρου.
Μονάχη έγνοια η γλώσσα μου στις αμμουδιές του Ομήρου.
από το Άξιον Εστί
του Οδυσσέα Ελύτη`

// Odysseus Elytis - Polytonic
var test_string_u6 = `Τὴ γλῶσσα μοῦ ἔδωσαν ἑλληνικὴ
τὸ σπίτι φτωχικὸ στὶς ἀμμουδιὲς τοῦ Ὁμήρου.
Μονάχη ἔγνοια ἡ γλῶσσα μου στὶς ἀμμουδιὲς τοῦ Ὁμήρου.
ἀπὸ τὸ Ἄξιον ἐστί
τοῦ Ὀδυσσέα Ἐλύτη`

// Pushkin's Bronze Horseman
var test_string_u7 = `На берегу пустынных волн
Стоял он, дум великих полн,
И вдаль глядел. Пред ним широко
Река неслася; бедный чёлн
По ней стремился одиноко.
По мшистым, топким берегам
Чернели избы здесь и там,
Приют убогого чухонца;
И лес, неведомый лучам
В тумане спрятанного солнца,
Кругом шумел.`

var list_test_string_u1_7 = []string{
	test_string_u1,
	test_string_u2,
	test_string_u3,
	test_string_u4,
	test_string_u5,
	test_string_u6,
	test_string_u7,
}

var test_sizes1 = []uint{0, 1, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192}
var test_runes1 = []rune(test_string1)
var test_runes_09 = []rune(test_string_09)
var test_runes_u1 = []rune(test_string_u1)
var test_runes_u2 = []rune(test_string_u2)
var test_runes_u3 = []rune(test_string_u3)
var test_runes_u4 = []rune(test_string_u4)
var test_runes_u5 = []rune(test_string_u5)
var test_runes_u6 = []rune(test_string_u6)
var test_runes_u7 = []rune(test_string_u7)

var list_test_runes_u1_7 = [][]rune{
	test_runes_u1,
	test_runes_u2,
	test_runes_u3,
	test_runes_u4,
	test_runes_u5,
	test_runes_u6,
	test_runes_u7,
}

func push_runes(rb *RuneBuff, runes []rune) {
	for _, r := range runes {
		rb.PushRune(r)
	}
}

func TestRuneBuffCapacity(t *testing.T) {
	tag := "TestRuneBuffCapacity"
	sizes := test_sizes1
	for _, size := range sizes {
		rb := NewRuneBuff(size)
		c := rb.Capacity()
		if c != size {
			// special case we always want capacity to be >= 1
			if !(c == KRUNE_BUFF_MIN && size < KRUNE_BUFF_MIN) {
				t.Fatalf(tag+" fail, expected %v, found %v", size, c)
			}
		}
	}
}

func TestRuneBuffReset(t *testing.T) {
	tag := "TestRuneBuffReset"
	sizes := test_sizes1
	for _, size := range sizes {
		rb := NewRuneBuff(size)
		rb.Reset()
		index := rb.Index()
		if index != 0 {
			t.Fatalf(tag+" fail, expected zero, found %v", index)
		}
	}
}

func TestRuneBuffContains(t *testing.T) {
	tag := "TestRuneBuffContains"

	for _, runes := range list_test_runes_u1_7 {
		rb := NewRuneBuff(0)
		push_runes(rb, runes)
		for index, r := range runes {
			if !rb.Contains(r) {
				t.Fatalf(tag+" failed, expected to find rune (%v)\nsample: %v\nindex: %v", r, runes, index)
			}
		}
	}
}

func TestRuneBuffFirst(t *testing.T) {
	tag := "TestRuneBuffFirst"
	for _, runes := range list_test_runes_u1_7 {
		for _, r := range runes {
			rb := NewRuneBuff(0)
			rb.PushRune(r)
			f := rb.First()
			if f != r {
				t.Fatalf(tag+" failed, expected rune %v, found %v", r, f)
			}
		}
	}
}

func TestRuneBuffIndexOf(t *testing.T) {
	tag := "TestRuneBuffIndexOf"
	runes := test_runes_09
	rb := NewRuneBuff(0)
	push_runes(rb, runes)
	for index, r := range runes {
		rbx := rb.IndexOf(r)
		if rbx != index {
			t.Fatalf(tag+" failed, expected index %v, found %v, rune %v", index, rbx, string(r))
		}
	}
}

func TestRuneBuffString(t *testing.T) {
	tag := "TestRuneBuffString"
	for _, runes := range list_test_runes_u1_7 {
		rb := NewRuneBuff(0)
		push_runes(rb, runes)
		s := rb.String()
		srunes := string(runes)
		if s != srunes {
			t.Fatalf(tag+" failed, expected [%v], found [%v]", srunes, s)
		}
	}
}

func TestRuneBuffPushRune(t *testing.T) {
	tag := "TestRuneBuffPushRune"
	for _, runes := range list_test_runes_u1_7 {
		rb := NewRuneBuff(0)
		rix := uint(0)
		for _, r := range runes {
			rix++
			rb.PushRune(r)
			if rb.Index() != rix {
				t.Fatalf(tag+" failed, expected index after push %v, found %v", rix, rb.Index())
			}
			rrb, err := rb.AtIndex(rix - 1)
			if err != nil {
				t.Fatalf(tag+" failed unexpected error: %v", err.Error())
			}
			if rrb != r {
				t.Fatalf(tag+" failed, doesn't agree with itself, expected rune %v, found %v", r, rrb)
			}
		}
	}
}

func TestRuneBuffPushString(t *testing.T) {
	tag := "TestRuneBuffPushString"

	for _, s := range list_test_string_u1_7 {
		rb := NewRuneBuff(0)
		rb.PushString(s)
		srb := rb.String()
		if srb != s {
			t.Fatalf(tag+" failed, expected string [%v], found [%v]", s, srb)
		}
	}
}

func TestRuneBuffSetIndex(t *testing.T) {
	tag := "TestRuneBuffSetIndex"

	for _, runes := range list_test_runes_u1_7 {
		rb := NewRuneBuff(0)
		push_runes(rb, runes)
		h := len(runes) / 2
		rb.SetIndex(uint(h))
		rbx := rb.Index()
		if rbx != uint(h) {
			t.Fatalf(tag+" failed, expected index %v, found index %v", h, rbx)
		}
		s := string(runes[:h])
		srb := rb.String()
		if s != srb {
			t.Fatalf(tag+" String() failed, expected [%v], found [%v]", s, srb)
		}
	}
}

func TestRuneBuffWriteTo(t *testing.T) {
	tag := "TestRuneBuffWriteTo"

	for _, runes := range list_test_runes_u1_7 {
		rb := NewRuneBuff(0)
		push_runes(rb, runes)
		bbuf := bytes.NewBuffer([]byte{})
		b, err := rb.WriteTo(bbuf)
		if err != nil {
			t.Fatalf(tag+" unexpected error on write: %v", err.Error())
		}
		rbx := rb.Index()
		if b < int(rbx) {
			// b should be >= number of runes in rb
			t.Fatalf(tag+" unexpected write size, should be at least %v, found %v", rbx, b)
		}
		sbbuf := bbuf.String()
		srunes := string(runes)
		if sbbuf != srunes {
			t.Fatalf(tag+" validation failed, expected [%v] found [%v]", srunes, sbbuf)
		}
	}
}
