/*
 * Copyright (c) 2025 Fumiama Minamoto.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package blake2b

import (
	"bytes"
	"crypto/rand"
	"encoding"
	"encoding/hex"
	"testing"
)

func TestBinaryMarshalUnmarshal256(t *testing.T) {
	buf := make([]byte, 16384)
	rand.Read(buf)
	var (
		exp [32]byte
		got [32]byte
	)
	for i := 0; i < 16384; i++ {
		h := New256()
		_, err := h.Write(buf[:i])
		if err != nil {
			t.Fatal(err)
		}
		data, err := h.(encoding.BinaryMarshaler).MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		newh := New256()
		err = newh.(encoding.BinaryUnmarshaler).UnmarshalBinary(data)
		if err != nil {
			t.Fatal(err)
		}
		_, err = h.Write(buf[i:])
		if err != nil {
			t.Fatal(err)
		}
		_, err = newh.Write(buf[i:])
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(h.Sum(exp[:0]), newh.Sum(got[:0])) {
			t.Fatal("Expect", hex.EncodeToString(exp[:]), "but got", hex.EncodeToString(got[:]))
		}
	}
}
