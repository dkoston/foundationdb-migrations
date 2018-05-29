package trace

import (
	"testing"
)

func something() string {
	return Trace()
}

func other() string {
	return nestedOnce()
}

func nestedOnce() string {
	return Trace()
}


func first() string {
	return second()
}

func second() string {
	return third()
}

func third() string {
	return Trace()
}


func TestTrace(t *testing.T) {
	var expected = "command-line-arguments.something"
	var called = something()

	if expected != called {
		t.Errorf("Expected %s, called %s", expected, called)
	}

	var oneNestExpected = "command-line-arguments.nestedOnce"
	var oneNestCalled = other()

	if oneNestExpected != oneNestCalled {
		t.Errorf("oneNestExpected %s, oneNestCalled %s", oneNestExpected, oneNestCalled)
	}

	var thirdExpected = "command-line-arguments.third"
	var thirdCalled = first()

	if thirdExpected != thirdCalled {
		t.Errorf("thirdExpected %s, thirdCalled %s", thirdExpected, thirdCalled)
	}
}