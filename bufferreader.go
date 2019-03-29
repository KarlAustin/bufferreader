package bufferreader
/**
 * Based on: https://github.com/pisdhooy/fmtbytes/blob/master/fmtbytes.go
 *
 * TODO: Raise errors instead of ignore
 */

import (
 "bytes"
 "encoding/binary"
)

type BufferReader struct {
 data *bytes.Reader
 position int64
}

func NewBufferReader(data []byte) *BufferReader {
 return &BufferReader{bytes.NewReader(data), 0}
}

func (br *BufferReader) ReadNInt(len int) []byte {
 br.position += int64(len)
 ret := make([]byte, len)
 br.data.Read(ret)
 return ret
}

func (br *BufferReader) ReadLongLong() uint64 {
 br.position += 8
 ret := make([]byte, 8)
 br.data.Read(ret)
 return binary.BigEndian.Uint64(ret)
}

func (br *BufferReader) ReadLong() uint32 {
 br.position += 4
 ret := make([]byte, 4)
 br.data.Read(ret)
 return binary.BigEndian.Uint32(ret)
}

func (br *BufferReader) ReadShort() uint16 {
 br.position += 2
 ret := make([]byte, 2)
 br.data.Read(ret)
 return binary.BigEndian.Uint16(ret)
}

func (br *BufferReader) ReadSingle() int {
 br.position += 1
 ret := make([]byte, 1)
 br.data.Read(ret)
 return int(ret[0])
}

func (br *BufferReader) ReadString(len int) string {
 br.position += int64(len)
 ret := make([]byte, len)
 br.data.Read(ret)
 return string(ret)
}

func (br *BufferReader) ReadArray16(len uint32) []uint16 {
 /**
  * We could do this with multiple calls to ReadShort() but it's really inefficient as it has to create a byte
  * slice for every call to ReadShort, here we just create 2 in total.  In the worst case on a single uint16 it's
  * the same total of slices.
  */
 br.position += int64(len)
 bufferSize := len / 2
 ret := make([]uint16, bufferSize)
 raw := br.ReadRaw(int(len))
 var i uint32
 var start = 0
 var val = make([]byte, 2)
 for i = 0; i < bufferSize; i++ {
  start = int(i) * 2
  val = raw[start:start+2]
  ret = append(ret, binary.BigEndian.Uint16(val))
 }
 return ret
}

func (br *BufferReader) ReadArray8(len uint32) []uint8 {
 /**
  * We could do this with multiple calls to ReadSingle() but it's really inefficient as it has to create a byte
  * slice for every call to ReadSingle, here we just create 1 in total.  In the worst case of a single uint8 it's
  * 1 slice vs. 2 for ReadSingle().
  */
 br.position += int64(len)
 bufferSize := len
 ret := make([]uint8, bufferSize)
 raw := br.ReadRaw(int(len))
 var i uint32
 var val byte
 for i = 0; i < bufferSize; i++ {
  val = raw[int(i)]
  ret = append(ret, uint8(val))
 }
 return ret
}

func (br *BufferReader) ReadRaw(len int) []byte {
 br.position += int64(len)
 ret := make([]byte, len)
 br.data.Read(ret)
 return ret
}

func (br *BufferReader) ReadRawOffset(len int, offset int64) []byte {
 br.position += int64(len)
 ret := make([]byte, len)
 br.data.ReadAt(ret, offset)
 return ret
}

func (br *BufferReader) ReadDouble() float64 {
 br.position += 8
 buffer := make([]byte, 8)
 br.data.Read(buffer)
 read := bytes.NewReader(buffer)
 var ret float64
 binary.Read(read, binary.BigEndian, &ret)
 return ret
}

func (br *BufferReader) ReadFloat() float32 {
 br.position += 4
 buffer := make([]byte, 4)
 br.data.Read(buffer)
 read := bytes.NewReader(buffer)
 var ret float32
 binary.Read(read, binary.BigEndian, &ret)
 return ret
}

/*func (br *BufferReader) ReadU16f16(littleEndian bool) (uint16, float32) {
 var mantissa float64
 var tExp uint16
 var tMan uint16

 exp := br.ReadRaw(2)
 man := br.ReadRaw(2)

 if littleEndian {
  tExp = binary.LittleEndian.Uint16(exp)
  tMan = binary.LittleEndian.Uint16(man)
 } else {
  tExp = binary.BigEndian.Uint16(exp)
  tMan = binary.BigEndian.Uint16(man)
 }

 bin := strconv.FormatUint(uint64(tMan), 2)

 len := len(bin)
 offset := 16 - len
 for i := 0; i < len; i++ {
  if string(bin[i]) == "1" {
   pow := (i+1+offset)*-1
   mantissa += 1 * math.Pow(2, float64(pow))
  }
 }

 return tExp, float32(mantissa)

}*/

func (br *BufferReader) Remaining() int64 {
 return br.data.Size() - br.Position()
}

func (br *BufferReader) Position() int64 {
 return br.position
}