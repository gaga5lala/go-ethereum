// Copyright 2017 AMIS Technologies
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package pbft

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	testAddress  = "0x70524d664ffe731100208a0154e556f9bb679ae6"
	testAddress2 = "0xb37866a925bccd69cfa98d43b510f1d23d78a851"
)

func TestValidatorSet(t *testing.T) {
	testNewValidatorSet(t)
	testNormalValSet(t)
	testEmptyValSet(t)
}

func testNewValidatorSet(t *testing.T) {
	var validators []*Validator
	const ValCnt = 100

	// Create 100 validators with random addresses
	for i := 0; i < ValCnt; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		val := NewValidator(addr)
		validators = append(validators, val)
	}

	// Create ValidatorSet
	validatorSet := NewValidatorSet(validators)

	// Check validators sorting: should be in ascending order
	for i := 0; i < ValCnt-1; i++ {
		val := validatorSet.GetByIndex(uint64(i))
		nextVal := validatorSet.GetByIndex(uint64(i + 1))
		if strings.Compare(val.Address().Hex(), nextVal.Address().Hex()) >= 0 {
			t.Errorf("Validator set is not sorted in sorted in ascending order")
		}
	}
}

func testNormalValSet(t *testing.T) {
	addr1 := common.HexToAddress(testAddress)
	val1 := NewValidator(addr1)
	addr2 := common.HexToAddress(testAddress2)
	val2 := NewValidator(addr2)

	valSet := NewValidatorSet([]*Validator{val1, val2})

	// check size
	if size := valSet.Size(); size != 2 {
		t.Errorf("wrong peer set size, got: %v, expected: 2", size)

	}
	// test get by index
	if val := valSet.GetByIndex(uint64(0)); val != val1 {
		t.Errorf("get wrong validator, got: %v, expected: %v", val, val1)
	}
	// test get by invalid index
	if val := valSet.GetByIndex(uint64(2)); val != nil {
		t.Errorf("get wrong validator, got: %v, expected: nil", val)
	}
	// test get by address
	if val := valSet.GetByAddress(addr2); val != val2 {
		t.Errorf("get wrong validator, got: %v, expected: %v", val, val2)
	}
	// test get by invalid address
	invalidAddr := common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
	if val := valSet.GetByAddress(invalidAddr); val != nil {
		t.Errorf("get wrong validator, got: %v, expected: nil", val)
	}
	// test get proposer
	if val := valSet.GetProposer(); val != val1 {
		t.Errorf("get wrong proposer, got: %v, expected: %v", val, val1)
	}
	// test calculate proposer
	valSet.CalcProposer(uint64(3))
	if val := valSet.GetProposer(); val != val2 {
		t.Errorf("get wrong proposer, got: %v, expected: %v", val, val2)
	}
}

func testEmptyValSet(t *testing.T) {
	valSet := NewValidatorSet([]*Validator{})

	// check size
	if size := valSet.Size(); size != 0 {
		t.Errorf("wrong peer set size, got: %v, expected: 0", size)

	}
	// test get by index
	if val := valSet.GetByIndex(uint64(0)); val != nil {
		t.Errorf("get wrong validator, got: %v, expected: nil", val)
	}
	// test get by invalid address
	invalidAddr := common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
	if val := valSet.GetByAddress(invalidAddr); val != nil {
		t.Errorf("get wrong validator, got: %v, expected: nil", val)
	}
	// test get proposer
	if val := valSet.GetProposer(); val != nil {
		t.Errorf("get wrong proposer, got: %v, expected: nil", val)
	}
	// test calculate proposer
	valSet.CalcProposer(uint64(3))
	if val := valSet.GetProposer(); val != nil {
		t.Errorf("get wrong proposer, got: %v, expected: nil", val)
	}
}
