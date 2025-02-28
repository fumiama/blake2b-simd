/*
 * Copyright (c) 2022 Fumiama Minamoto.
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

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"

	blake2b "github.com/fumiama/blake2b-simd"
)

func main() {
	help := flag.Bool("h", false, "display this help")
	size := flag.Int("s", 32, "digest output bytes size")
	key := flag.String("k", "", "hex-encoded key for prefix-MAC")
	var keyb []byte
	var err error
	if *key != "" {
		keyb, err = hex.DecodeString(*key)
		if err != nil {
			panic(err)
		}
	}
	salt := flag.String("t", "", "hex-encoded 16 bytes salt (if < 16 bytes, padded with zeros)")
	var saltb []byte
	if *salt != "" {
		saltb, err = hex.DecodeString(*salt)
		if err != nil {
			panic(err)
		}
	}
	pers := flag.String("p", "", "hex-encoded 16 bytes personalization (if < 16 bytes, padded with zeros)")
	var persb []byte
	if *pers != "" {
		persb, err = hex.DecodeString(*pers)
		if err != nil {
			panic(err)
		}
	}
	flag.Parse()
	if *help {
		fmt.Println(os.Args[0], "[commands] [file1 file2...(null for stdin)]")
		flag.PrintDefaults()
		return
	}
	h, err := blake2b.New(&blake2b.Config{
		Size:   uint8(*size),
		Key:    keyb,
		Salt:   saltb,
		Person: persb,
	})
	if err != nil {
		panic(err)
	}
	file := flag.Args()
	var f []io.Reader
	if len(file) == 0 {
		f = append(f, os.Stdin)
		file = append(file, "stdin")
	} else {
		for _, name := range file {
			fi, err := os.Open(name)
			if err != nil {
				panic(err)
			}
			f = append(f, fi)
		}
	}
	for i, fi := range f {
		_, err = io.Copy(h, fi)
		if err != nil {
			panic(err)
		}
		sum := make([]byte, 0, *size)
		fmt.Println(file[i], hex.EncodeToString(h.Sum(sum)))
		h.Reset()
	}
}
