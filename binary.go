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
	"encoding/binary"
	"unsafe"
)

// MarshalBinary marshals current digest status into bytes.
func (d *digest) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, unsafe.Sizeof(*d)))
	_ = binary.Write(buf, binary.BigEndian, &d.h)
	_ = binary.Write(buf, binary.BigEndian, &d.t)
	_ = binary.Write(buf, binary.BigEndian, &d.f)
	_ = binary.Write(buf, binary.BigEndian, &d.x)
	_ = binary.Write(buf, binary.BigEndian, int64(d.nx))
	_ = binary.Write(buf, binary.BigEndian, &d.ih)
	_ = binary.Write(buf, binary.BigEndian, &d.paddedKey)
	if d.isKeyed {
		buf.WriteByte(1)
	} else {
		buf.WriteByte(0)
	}
	buf.WriteByte(d.size)
	if d.isLastNode {
		buf.WriteByte(1)
	} else {
		buf.WriteByte(0)
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary loads current digest status from bytes.
func (d *digest) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.BigEndian, &d.h)
	if err != nil {
		return err
	}
	err = binary.Read(r, binary.BigEndian, &d.t)
	if err != nil {
		return err
	}
	err = binary.Read(r, binary.BigEndian, &d.f)
	if err != nil {
		return err
	}
	err = binary.Read(r, binary.BigEndian, &d.x)
	if err != nil {
		return err
	}
	nx := int64(0)
	err = binary.Read(r, binary.BigEndian, &nx)
	if err != nil {
		return err
	}
	d.nx = int(nx)
	err = binary.Read(r, binary.BigEndian, &d.ih)
	if err != nil {
		return err
	}
	err = binary.Read(r, binary.BigEndian, &d.paddedKey)
	if err != nil {
		return err
	}
	b, err := r.ReadByte()
	if err != nil {
		return err
	}
	d.isKeyed = b != 0
	d.size, err = r.ReadByte()
	if err != nil {
		return err
	}
	b, err = r.ReadByte()
	if err != nil {
		return err
	}
	d.isLastNode = b != 0
	return nil
}
