package e2etests

import "testing"

var emptyStat = stats{0, 0, 0}

func TestEmptyStats(t *testing.T) {
	err := registerNewUser("statsUser", "test")
	if err != nil {
		t.Error(err)
	}

	stats, err := awaitStats("statsUser", "test")
	if err != nil {
		t.Error(err)
	}

	assertStatsEmpty(stats, emptyStat, t)
}

func assertStatsEmpty(got, want stats, t *testing.T) {
	if got != want {
		t.Errorf("stats are not equal, got %v but want %v", got, want)
	}
}
