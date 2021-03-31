/*
 * Copyright (c) 2014 Mark Samman <https://github.com/marksamman/bencode>
 * Copyright (c) 2021 Aleksandr Panov <https://github.com/impos/bencode>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package bencode

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

const (
	typeDictionary = 'd'
	typeList       = 'l'
	typeInteger    = 'i'
	endOfValue     = 'e'
)

type encoder struct {
	data *bytes.Buffer
}

func (e *encoder) writeInt(i int64) error {
	e.data.WriteByte(typeInteger)
	e.data.WriteString(strconv.FormatInt(i, 10))
	e.data.WriteByte(endOfValue)

	return nil
}

func (e *encoder) writeUint(i uint64) error {
	e.data.WriteByte(typeInteger)
	e.data.WriteString(strconv.FormatUint(i, 10))
	e.data.WriteByte(endOfValue)

	return nil
}

func (e *encoder) writeString(s string) error {
	e.data.WriteString(strconv.Itoa(len(s)))
	e.data.WriteString(":")
	e.data.WriteString(s)

	return nil
}

func (e *encoder) writeList(l []interface{}) error {
	e.data.WriteByte(typeList)

	for _, v := range l {
		err := e.writeValue(v)
		if err != nil {
			return err
		}
	}

	e.data.WriteByte(endOfValue)

	return nil
}

func (e *encoder) writeDictionary(d map[string]interface{}) error {
	keys := make(sort.StringSlice, len(d))

	i := 0

	for k := range d {
		keys[i] = k
		i++
	}

	keys.Sort()

	e.data.WriteByte(typeDictionary)

	for _, k := range keys {
		if err := e.writeString(k); err != nil {
			return err
		} else if err = e.writeValue(d[k]); err != nil {
			return err
		}
	}

	e.data.WriteByte(endOfValue)

	return nil
}

func (e *encoder) writeValue(i interface{}) error {
	switch v := i.(type) {
	case string:
		return e.writeString(v)
	case []byte:
		return e.writeString(string(v))
	case int, int8, int16, int32, int64:
		return e.writeInt(reflect.ValueOf(v).Int())
	case uint, uint16, uint32, uint64:
		return e.writeUint(reflect.ValueOf(v).Uint())
	case []interface{}:
		return e.writeList(v)
	case map[string]interface{}:
		return e.writeDictionary(v)
	default:
		return fmt.Errorf("unsupported value type: %T", i)
	}
}

func Encode(i map[string]interface{}) ([]byte, error) {
	e := encoder{bytes.NewBuffer(make([]byte, 0))}

	if err := e.writeDictionary(i); err != nil {
		return nil, fmt.Errorf("write value: %w", err)
	}

	return e.data.Bytes(), nil
}
