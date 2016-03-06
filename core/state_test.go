package core

import "testing"

func TestStateCopy(t *testing.T) {
	st1 := NewState()
	st1.Uniform("test").Set(0)

	st2 := st1
	st2.Uniform("test").Set(1)

	if st1.Uniform("test").Value().(int) != 1 {
		// we expect map uniform to be shared, so contaminated
		t.Fail()
	}

	// make sure we got the same uniforms
	st3 := st1.Copy()
	if st3.Uniform("test").Value().(int) != 1 {
		t.Fatal("st3 unexpected uniform value")
	}

	// make sure we can change st3 without messing with st1
	st3.Uniform("test").Set(2)

	if st1.Uniform("test").Value().(int) == 2 {
		t.Fatal("st1 value contaminated by st3 assignment")
	}
	t.Log(st1.Uniform("test").Value())
}
