package core

import "testing"

func TestConditions(t *testing.T) {
	tables := []struct {
		cond Condition
		F uint8
		res bool
	}{
		{NonZero   , 0, true},
		{NonZero   , FlagZ, false},
		{Zero      , FlagZ, true},
		{Zero      , 0, false},
		
		{NoCarry   , 0, true},
		{NoCarry   , FlagC, false},
		{Carry     , FlagC, true},
		{Carry     , 0, false},

		{ParityOdd , 0, true},
		{ParityOdd , FlagP, false},
		{ParityEven, FlagP, true},
		{ParityEven, 0, false},

		{SignPos   , 0, true},
		{SignPos   , FlagS, false},
		{SignNeg   , FlagS, true},
		{SignNeg   , 0, false},
	}

	for _, table := range tables {
		res := table.cond.isTrue(&table.F)
		
		if res != table.res {
			t.Errorf("Condition %v (%#04x) with Flag bits %#02x was incorrectly evaluated, got: %v, want: %v.", table.cond, table.cond, table.F , res, table.res)
		}
	}
}