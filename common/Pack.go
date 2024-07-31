package common

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Pack struct {
	method   int32
	err      int32
	sequence int32
	data     []byte
	position int
}

func NewPack() *Pack {
	return &Pack{}
}

func NewPackInstance(byteBuffer []byte) *Pack {
	buffer := bytes.NewReader(byteBuffer)
	instance := NewPack()

	binary.Read(buffer, binary.BigEndian, &instance.method)
	binary.Read(buffer, binary.BigEndian, &instance.err)
	binary.Read(buffer, binary.BigEndian, &instance.sequence)

	instance.data = make([]byte, len(byteBuffer)-12)
	buffer.Read(instance.data)

	return instance
}

func NewResponse(pack *Pack) *Pack {
	response := NewPack()
	response.SetMethod(pack.method)
	response.SetSequence(pack.sequence)
	return response
}

func (p *Pack) ToByteBuff() []byte {
	size := int32(len(p.data) + 12)
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, size)
	binary.Write(buffer, binary.BigEndian, p.method)
	binary.Write(buffer, binary.BigEndian, p.err)
	binary.Write(buffer, binary.BigEndian, p.sequence)
	buffer.Write(p.data)
	return buffer.Bytes()
}

func (p *Pack) ReadInt() int32 {
	if !p.IsHasInt() {
		return 0
	}
	var result int32
	binary.Read(bytes.NewReader(p.data[p.position:p.position+4]), binary.BigEndian, &result)
	p.position += 4
	return result
}

func (p *Pack) IsHasInt() bool {
	return len(p.data)-p.position >= 4
}

func (p *Pack) ReadFloat() float32 {
	if !p.IsHasFloat() {
		return 0
	}
	var result float32
	binary.Read(bytes.NewReader(p.data[p.position:p.position+4]), binary.BigEndian, &result)
	p.position += 4
	return result
}

func (p *Pack) IsHasFloat() bool {
	return p.IsHasInt()
}

func (p *Pack) ReadLong() int64 {
	if !p.IsHasLong() {
		return 0
	}
	var result int64
	binary.Read(bytes.NewReader(p.data[p.position:p.position+8]), binary.BigEndian, &result)
	p.position += 8
	return result
}

func (p *Pack) IsHasLong() bool {
	return len(p.data)-p.position >= 8
}

func (p *Pack) ReadBytes() []byte {
	length := p.ReadInt()
	result := make([]byte, length)
	copy(result, p.data[p.position:p.position+int(length)])
	p.position += int(length)
	return result
}

func (p *Pack) IsHasBytes() bool {
	return p.IsHasInt()
}

func (p *Pack) ReadString() string {
	return string(p.ReadBytes())
}

func (p *Pack) IsReadable() bool {
	return len(p.data)-p.position > 0
}

func (p *Pack) WriteInt(value int32) {
	p.checkBuffer(4)
	binary.BigEndian.PutUint32(p.data[p.position:], uint32(value))
	p.position += 4
}

func (p *Pack) WriteFloat(value float32) {
	p.checkBuffer(4)
	binary.BigEndian.PutUint32(p.data[p.position:], math.Float32bits(value))
	p.position += 4
}

func (p *Pack) WriteLong(value int64) {
	p.checkBuffer(8)
	binary.BigEndian.PutUint64(p.data[p.position:], uint64(value))
	p.position += 8
}

func (p *Pack) WriteBytes(bytes []byte) {
	p.WriteInt(int32(len(bytes)))
	p.checkBuffer(len(bytes))
	copy(p.data[p.position:], bytes)
	p.position += len(bytes)
}

func (p *Pack) WriteString(s string) {
	p.WriteBytes([]byte(s))
}

func (p *Pack) checkBuffer(size int) {
	if len(p.data)-p.position < size {
		newData := make([]byte, len(p.data)+size)
		copy(newData, p.data)
		p.data = newData
	}
}

func (p *Pack) WriteProto(message protoreflect.ProtoMessage) {
	bytes, err := proto.Marshal(message)
	if err != nil {
		log.Fatalln("Error marshalling proto:", err)
		return
	}
	p.WriteBytes(bytes)
}

func (p *Pack) ReadProto(message protoreflect.ProtoMessage) {
	if err := proto.Unmarshal(p.ReadBytes(), message); err != nil {
		log.Fatalln("Failed to parse proto message:", err)
	}
}

func (p *Pack) SetSequence(value int32) {
	p.sequence = value
}

func (p *Pack) GetSequence() int32 {
	return p.sequence
}

func (p *Pack) SetMethod(value int32) {
	p.method = value
}

func (p *Pack) GetMethod() int32 {
	return p.method
}

func (p *Pack) SetError(value bool) {
	if value {
		p.err = 1
	} else {
		p.err = 0
	}
}

func (p *Pack) GetError() bool {
	return p.err > 0
}

func (p *Pack) GetData() []byte {
	return p.data
}
