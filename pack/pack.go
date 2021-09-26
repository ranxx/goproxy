package pack

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

// Package ...
type Package struct {
	Version   [2]byte // 协议版本
	Length    int64   // 数据部分长度
	Timestamp int64   // 时间戳
	Msg       []byte  // 数据部分长度
}

// PackBytes ...
func (p *Package) PackBytes() ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 2+8+8+len(p.Msg)))
	err := p.Pack(buffer)
	return buffer.Bytes(), err
}

// Pack ...
func (p *Package) Pack(writer io.Writer) error {
	var err error
	err = binary.Write(writer, binary.BigEndian, &p.Version)
	err = binary.Write(writer, binary.BigEndian, &p.Length)
	err = binary.Write(writer, binary.BigEndian, &p.Timestamp)
	err = binary.Write(writer, binary.BigEndian, &p.Msg)
	return err
}

// UnpackBytes ...
func (p *Package) UnpackBytes(body []byte) error {
	buffer := bytes.NewBuffer(body)
	return p.Unpack(buffer)
}

// Unpack ...
func (p *Package) Unpack(reader io.Reader) error {
	var err error
	err = binary.Read(reader, binary.BigEndian, &p.Version)
	err = binary.Read(reader, binary.BigEndian, &p.Length)
	err = binary.Read(reader, binary.BigEndian, &p.Timestamp)
	p.Msg = make([]byte, p.Length-2-8-8)
	err = binary.Read(reader, binary.BigEndian, &p.Msg)
	return err
}

// SplitFunc split
func SplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	fmt.Printf("%t|%d|%s\n", atEOF, len(data), data)
	if len(data) > 2 {
		fmt.Println(atEOF, string(data[0]), string(data[1]))
	}
	// for i, v := range data {
	// 	fmt.Printf("%d|%c, ", i, v)
	// }
	if !atEOF && data[0] == 'X' && data[1] == '1' {
		if len(data) > 10 {
			// fmt.Println("开始解析")
			length := int64(0)
			binary.Read(bytes.NewReader(data[2:10]), binary.BigEndian, &length)

			// fmt.Println("解析长度为", length, len(data))

			if int(length) <= len(data) {
				// fmt.Println("开始解析-2")
				return int(length), data[:int(length)], nil
			}
		}
	}
	return
}

// NewPackage ...
func NewPackage(msg []byte) *Package {
	pack := Package{
		Version:   [2]byte{'X', '1'},
		Length:    int64(len(msg) + 2 + 8 + 8),
		Timestamp: time.Now().Unix(),
		Msg:       msg,
	}
	return &pack
}

// NewScanner ...
func NewScanner(reader io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(reader)
	scanner.Split(SplitFunc)
	return scanner
}
