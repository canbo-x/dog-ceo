package dummy_rate_limiter

import (
	"reflect"
	"testing"
	"time"
)

func TestNewLimitCounter(t *testing.T) {
	lc := NewLimitCounter()

	if reflect.TypeOf(lc) != reflect.TypeOf(&LimitCounter{}) {
		t.Errorf("NewLimitCounter method does not return imitCounter instance correctly")
	}

}

func TestLimitCounter(t *testing.T) {
	lc := NewLimitCounter()

	count := lc.get()
	if count != 0 {
		t.Errorf("get supposed to return 0 but returned %d", count)
	}

	lc.increase()
	count++

	if c := lc.get(); c != count {
		t.Errorf("get supposed to return %d but returned %d", count, c)
	}

	lc.decrease()
	count--
	if c := lc.get(); c != count {
		t.Errorf("get supposed to return %d but returned %d", count, c)
	}

	for i := 0; i < 10; i++ {
		count++
		lc.increase()
	}

	if !lc.isLimitReached() {
		t.Errorf("isLimitReached supposed to return true but returned false")
	}

	lc.decrease()
	count--
	if lc.isLimitReached() {
		t.Errorf("isLimitReached supposed to return false but returned true")
	}

	lc.StartLimiter()
	time.Sleep(time.Second * 10)
	if c := lc.get(); c != 0 {
		t.Errorf("get supposed to return 0 but returned %d", c)
	}

	for i := 0; i < 10; i++ {
		lc.decrease()
	}

	if c := lc.get(); c != 0 {
		t.Errorf("get supposed to return 0 but returned %d", c)
	}

}
