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
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type decoder struct {
	data *bufio.Reader
}

func (d *decoder) readValue(valType byte) (interface{}, error) {
	switch valType {
	case typeDictionary:
		return d.readDictionary()
	case typeList:
		return d.readList()
	case typeInteger:
		return d.readInteger()
	default:
		err := d.data.UnreadByte()
		if err != nil {
			return nil, fmt.Errorf("unread byte: %w", err)
		}

		return d.readString()
	}
}

func (d *decoder) readInteger() (interface{}, error) {
	valBytes, err := d.data.ReadBytes(endOfValue)
	if err != nil {
		return 0, fmt.Errorf("read bytes: %w", err)
	}

	valStr := string(valBytes[:len(valBytes)-1])

	if val, err := strconv.ParseInt(valStr, 10, 64); err == nil {
		return val, nil
	} else if val, err := strconv.ParseUint(valStr, 10, 64); err == nil {
		return val, nil
	} else {
		return 0, err
	}
}

func (d *decoder) readDictionary() (map[string]interface{}, error) {
	dict := map[string]interface{}{}

	for {
		key, err := d.readString()
		if err != nil {
			return dict, fmt.Errorf("read key: %w", err)
		}

		valType, err := d.data.ReadByte()
		if err != nil {
			return dict, fmt.Errorf("read value type: %w", err)
		}

		if valType == endOfValue {
			break
		}

		val, err := d.readValue(valType)
		if err != nil {
			return dict, fmt.Errorf("read value: %w", err)
		}

		dict[key] = val

		valType, err = d.data.ReadByte()
		if err != nil {
			return dict, fmt.Errorf("read dictionary end: %w", err)
		}

		if valType == endOfValue {
			break
		}

		err = d.data.UnreadByte()
		if err != nil {
			return dict, fmt.Errorf("unread dictionary end: %w", err)
		}
	}

	return dict, nil
}

func (d *decoder) readList() ([]interface{}, error) {
	list := make([]interface{}, 0)

	for {
		valType, err := d.data.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("read type: %w", err)
		}

		if valType == endOfValue {
			break
		}

		val, err := d.readValue(valType)
		if err != nil {
			return nil, fmt.Errorf("read value: %w", err)
		}

		list = append(list, val)
	}

	return list, nil
}

func (d *decoder) readString() (string, error) {
	length, err := d.readLength()
	if err != nil {
		return "", fmt.Errorf("read length: %w", err)
	}

	result := make([]byte, length)

	_, err = io.ReadFull(d.data, result)
	if err != nil {
		return "", fmt.Errorf("read value: %w", err)
	}

	return string(result), nil
}

func (d *decoder) readLength() (int64, error) {
	valBytes, err := d.data.ReadBytes(':')
	if err != nil {
		return 0, fmt.Errorf("read bytes: %w", err)
	}

	valStr := string(valBytes[:len(valBytes)-1])

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse int: %w", err)
	}

	return val, nil
}

func Decode(data io.Reader) (map[string]interface{}, error) {
	d := decoder{bufio.NewReader(data)}

	if rootType, err := d.data.ReadByte(); err != nil {
		return map[string]interface{}{}, fmt.Errorf("read root type: %w", err)
	} else if rootType != typeDictionary {
		return nil, errors.New("is not valid bencode file")
	}

	return d.readDictionary()
}
