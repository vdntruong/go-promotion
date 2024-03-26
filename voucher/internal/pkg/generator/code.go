package generator

import "math/rand"

const (
	defaultCodeLength uint = 8
	defaultCodePool   uint = 5
	defaultCharset         = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type CodeValidator interface {
	IsUnique(string) bool
}

type VoucherGenerator struct {
	validator  CodeValidator
	codeChan   chan string
	stopChan   chan struct{}
	pool       uint
	codeLength uint
}

type VoucherGeneratorOption func(*VoucherGenerator)

func VoucherGeneratorCodePool(pool uint) VoucherGeneratorOption {
	return func(vg *VoucherGenerator) { vg.pool = pool }
}

func VoucherGeneratorCodeLength(length uint) VoucherGeneratorOption {
	return func(vg *VoucherGenerator) { vg.codeLength = length }
}

func NewVoucherGenerator(validator CodeValidator, opt ...VoucherGeneratorOption) *VoucherGenerator {
	var vg = &VoucherGenerator{
		validator:  validator,
		codeChan:   make(chan string),
		stopChan:   make(chan struct{}),
		pool:       defaultCodePool,
		codeLength: defaultCodeLength,
	}
	for _, o := range opt {
		o(vg)
	}
	for i := 0; i < int(vg.pool); i++ {
		go vg.generateWorker()
	}
	return vg
}

func (vg *VoucherGenerator) generateWorker() {
	for {
		select {
		case <-vg.stopChan:
			return
		default:
			code := make([]byte, vg.codeLength)
			for i := range code {
				code[i] = defaultCharset[rand.Intn(len(defaultCharset))]
			}
			if vg.validator.IsUnique(string(code)) {
				vg.codeChan <- string(code)
			}
		}
	}
}

func (vg *VoucherGenerator) Stop() {
	close(vg.stopChan)
}

func (vg *VoucherGenerator) GetUniqueCode() string {
	code := <-vg.codeChan
	return code
}
