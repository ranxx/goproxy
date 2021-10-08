package pack

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"time"
)

// Package ...
type Package struct {
	Version   [2]byte // 协议版本
	Length    int64   // 数据部分长度
	Timestamp int64   // 时间戳秒
	Msg       []byte  // 数据部分长度
}

// Empty 空pack
func Empty() *Package {
	return NewPackage(nil)
}

// Reset 重置
func (p *Package) Reset(msg []byte) *Package {
	p.Timestamp = time.Now().Unix()
	p.Msg = msg
	return p
}

// PreLength 前置长度
func (p *Package) PreLength() int64 {
	return 18
}

// PackBytes ...
func (p *Package) PackBytes() ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 2+8+8+len(p.Msg)))
	err := p.Pack(buffer)
	return buffer.Bytes(), err
}

// ReadPre ...
func (p *Package) ReadPre(reader io.Reader) error {
	var err error
	if err = binary.Read(reader, binary.BigEndian, &p.Version); err != nil {
		return err
	}
	if err = binary.Read(reader, binary.BigEndian, &p.Length); err != nil {
		return err
	}
	if err = binary.Read(reader, binary.BigEndian, &p.Timestamp); err != nil {
		return err
	}
	return err
}

// ReadLast ...
func (p *Package) ReadLast(reader io.Reader) error {
	var err error
	p.Msg = make([]byte, p.Length-2-8-8)
	err = binary.Read(reader, binary.BigEndian, &p.Msg)
	return err
}

func absInt64(in int64) int64 {
	if in > 0 {
		return in
	}
	return in * -1
}

// IsPackage 是否为package数据
func (p *Package) IsPackage(inerval time.Duration) bool {
	tmp := Empty()
	if p.Version != tmp.Version {
		return false
	}

	in := time.Now().Unix() - p.Timestamp

	if absInt64(in) > int64(inerval/time.Second) {
		return false
	}
	return true
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
	log.Println(atEOF, len(data))
	pkg := Package{}
	if atEOF || len(data) < int(pkg.PreLength()) {
		return
	}

	if err = pkg.ReadPre(bytes.NewReader(data[:pkg.PreLength()])); err != nil {
		return
	}

	if pkg.IsPackage(time.Second * 10); err != nil {
		return
	}

	if pkg.Length > int64(len(data)) {
		return
	}
	log.Println(atEOF, len(data), pkg.Length)
	return int(pkg.Length), data[:pkg.Length], nil
}

// SplitFuncEdgeTriggered ...
func SplitFuncEdgeTriggered(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// log.Println("client.tcp", atEOF, len(data))
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if !atEOF {
		return len(data), data[:], nil
	}
	return 0, nil, nil
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
func NewScanner(reader io.Reader, split bufio.SplitFunc) *bufio.Scanner {
	scanner := bufio.NewScanner(reader)
	scanner.Split(split)
	scanner.Buffer(nil, 1024*1024*1024)
	return scanner
}
