package dhcp4

import (
	"encoding/binary"
	"net"
	"time"
)

// SelectOrderOrAll has same functionality as SelectOrder, except if the order
// param is nil, whereby all options are added (in arbitary order).
func (o Options) SelectOrderOrAll(order []byte) []Option {
	if order == nil {
		opts := make([]Option, 0, len(o))
		for i, v := range o {
			opts = append(opts, Option{Code: i, Value: v})
		}
		return opts
	}
	return o.SelectOrder(order)
}

// SelectOrder returns a slice of options ordered and selected by a byte array
// usually defined by OptionParameterRequestList.  This result is expected to be
// used in ReplyPacket()'s []Option parameter.
func (o Options) SelectOrder(order []byte) []Option {
	opts := make([]Option, 0, len(order))
	for _, v := range order {
		if data, ok := o[OptionCode(v)]; ok {
			opts = append(opts, Option{Code: OptionCode(v), Value: data})
		}
	}
	return opts
}

// IPRange returns how many ips in the ip range from start to stop (inclusive)
func IPRange(start, stop net.IP) int {
	//return int(Uint([]byte(stop))-Uint([]byte(start))) + 1
	return int(binary.BigEndian.Uint32(stop.To4())) - int(binary.BigEndian.Uint32(start.To4())) + 1
}

// IPAdd returns a copy of start + add.
// IPAdd(net.IP{192,168,1,1},30) returns net.IP{192.168.1.31}
func IPAdd(start net.IP, add int) net.IP { // IPv4 only
	start = start.To4()
	//v := Uvarint([]byte(start))
	result := make(net.IP, 4)
	binary.BigEndian.PutUint32(result, binary.BigEndian.Uint32(start)+uint32(add))
	//PutUint([]byte(result), v+uint64(add))
	return result
}

// IPLess returns where IP a is less than IP b.
func IPLess(a, b net.IP) bool {
	b = b.To4()
	for i, ai := range a.To4() {
		if ai != b[i] {
			return ai < b[i]
		}
	}
	return false
}

// IPInRange returns true if ip is between (inclusive) start and stop.
func IPInRange(start, stop, ip net.IP) bool {
	return !(IPLess(ip, start) || IPLess(stop, ip))
}

// OptionsLeaseTime - converts a time.Duration to a 4 byte slice, compatible
// with OptionIPAddressLeaseTime.
func OptionsLeaseTime(d time.Duration) []byte {
	leaseBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(leaseBytes, uint32(d/time.Second))
	//PutUvarint(leaseBytes, uint64(d/time.Second))
	return leaseBytes
}

// JoinIPs returns a byte slice of IP addresses, one immediately after the other
// This may be useful for creating multiple IP options such as OptionRouter.
func JoinIPs(ips []net.IP) (b []byte) {
	for _, v := range ips {
		b = append(b, v.To4()...)
	}
	return
}

/*
// PutUint writes value to a byte slice.
func PutUint(data []byte, value uint64) {
	for i := len(data) - 1; i >= 0; i-- {
		data[i] = byte(value % 256)
		value /= 256
	}
}

// Uint returns a value from a byte slice.
// Values requiring more than 64bits, won't work correctly
func Uint(data []byte) (ans uint64) {
	for _, b := range data {
		ans <<= 8
		ans += uint64(b)
	}
	return
}*/
