package fdbm

import (
	"testing"
)

func TestMigrationMapSortUp(t *testing.T) {

	ms := migrationSorter{}

	// insert in any order
	ms = append(ms, newMigration(20120000, "test", "./20120000_testing1.go"))
	ms = append(ms, newMigration(20128000, "test", "./20128000_testing8.go"))
	ms = append(ms, newMigration(20129000, "test", "./20129000_testing9.go"))
	ms = append(ms, newMigration(20127000, "test", "./20127000_testing7.go"))

	ms.Sort(true) // sort Upwards

	sorted := []int64{20120000, 20127000, 20128000, 20129000}

	validateMigrationSort(t, ms, sorted)
}

func TestMigrationMapSortDown(t *testing.T) {

	ms := migrationSorter{}

	// insert in any order
	ms = append(ms, newMigration(20120000, "test", "./20120000_testing1.go"))
	ms = append(ms, newMigration(20128000, "test", "./20128000_testing8.go"))
	ms = append(ms, newMigration(20129000, "test", "./20129000_testing9.go"))
	ms = append(ms, newMigration(20127000, "test", "./20127000_testing7.go"))

	ms.Sort(false) // sort Downwards

	sorted := []int64{20129000, 20128000, 20127000, 20120000}

	validateMigrationSort(t, ms, sorted)
}

func validateMigrationSort(t *testing.T, ms migrationSorter, sorted []int64) {

	for i, m := range ms {
		if sorted[i] != m.Number {
			t.Error("incorrect sorted number")
		}

		var next, prev int64

		if i == 0 {
			prev = -1
			next = ms[i+1].Number
		} else if i == len(ms)-1 {
			prev = ms[i-1].Number
			next = -1
		} else {
			prev = ms[i-1].Number
			next = ms[i+1].Number
		}

		if m.Next != next {
			t.Errorf("mismatched Next. v: %v, got %v, wanted %v\n", m, m.Next, next)
		}

		if m.Previous != prev {
			t.Errorf("mismatched Previous v: %v, got %v, wanted %v\n", m, m.Previous, prev)
		}
	}

	t.Log(ms)
}
