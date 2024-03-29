package vm

import (
	"math"
	"math/big"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/cloudwego/kitex/pkg/klog"
)

// Global Debug flag indicating Debug VM (full logging)
var Debug bool

type Type byte

const (
	StdVmTy Type = iota
	JitVmTy
	MaxVmTy

	LogTyPretty byte = 0x1
	LogTyDiff   byte = 0x2
)

var (
	Pow256 = common.BigPow(2, 256)

	U256 = common.U256
	S256 = common.S256

	Zero = common.Big0
	One  = common.Big1

	max = big.NewInt(math.MaxInt64)
)

func NewVm(env Environment) VirtualMachine {
	switch env.VmType() {
	case JitVmTy:
		return NewJitVm(env)
	default:
		klog.Error("unsupported vm type %d", env.VmType())
		fallthrough
	case StdVmTy:
		return New(env)
	}
}

func calcMemSize(off, l *big.Int) *big.Int {
	if l.Cmp(common.Big0) == 0 {
		return common.Big0
	}

	return new(big.Int).Add(off, l)
}

// Simple helper
func u256(n int64) *big.Int {
	return big.NewInt(n)
}

// Mainly used for print variables and passing to Print*
func toValue(val *big.Int) interface{} {
	// Let's assume a string on right padded zero's
	b := val.Bytes()
	if b[0] != 0 && b[len(b)-1] == 0x0 && b[len(b)-2] == 0x0 {
		return string(b)
	}

	return val
}

func getData(data []byte, start, size *big.Int) []byte {
	dlen := big.NewInt(int64(len(data)))

	s := common.BigMin(start, dlen)
	e := common.BigMin(new(big.Int).Add(s, size), dlen)
	return common.RightPadBytes(data[s.Uint64():e.Uint64()], int(size.Uint64()))
}

func UseGas(gas, amount *big.Int) bool {
	if gas.Cmp(amount) < 0 {
		return false
	}

	// Sub the amount of gas from the remaining
	gas.Sub(gas, amount)
	return true
}
