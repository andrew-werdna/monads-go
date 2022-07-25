package maybe

import "testing"

func TestMaybe(t *testing.T) {
	const testName1 = "test name"

	t.Run("should be able to Map between same type", func(t *testing.T) {
		type customType1 struct {
			Name string
			Age  int
		}

		myCustomStruct1 := customType1{
			Name: "test name",
			Age:  38,
		}
		m1 := Just[customType1, customType1]{
			Start: &myCustomStruct1,
		}

		m2 := m1.Map(func(t *customType1) *customType1 {
			t.Name = "new test name"
			return t
		})

		if m1 == m2 {
			t.Errorf("expected different address, got m1: %v, m2: %v\n", m1, m2)
		}

		if _, ok := m2.(Just[customType1, customType1]); !ok {
			t.Errorf("expected ok type assert, instead got: %v\n", ok)
		}
	})

	t.Run("should avoid nil pointer issues for same type", func(t *testing.T) {
		type customType1 struct {
			Name string
			Age  int
		}

		myCustomStruct1 := customType1{
			Name: "test name",
			Age:  38,
		}

		m1 := Just[customType1, customType1]{
			Start: &myCustomStruct1,
		}

		// m2 has a nil pointer in it
		m2 := m1.Map(func(t *customType1) *customType1 {
			var next *customType1
			return next
		}).Map(func(t *customType1) *customType1 {
			// we are protected against nil pointer dereference
			return &customType1{
				Name: "something",
				Age:  1,
			}
		})

		if n1, ok := m2.(Just[customType1, customType1]); ok {
			t.Fatalf("n1 should be Nothing, instead got %T", n1)
		}

		if n1, ok := m2.(Nothing[customType1, customType1]); !ok {
			t.Fatalf("n1 should be Nothing, instead got %T", n1)
		}
	})

	t.Run("should be able to map once between different types", func(t *testing.T) {
		type t1 struct {
			Name string
			Age  int
		}
		type t2 struct {
			NameLength int
		}

		vt1 := t1{
			Name: "test name",
			Age:  39,
		}

		m1 := Just[t1, t2]{
			Start: &vt1,
		}

		m2 := m1.Map(func(t *t1) *t2 {
			return &t2{
				NameLength: len(t.Name),
			}
		})

		if j2, ok := m2.(Just[t1, t2]); !ok {
			t.Fatalf("should be able to type assert ot Just[customType1, customType2]")
		} else {
			var _ Maybe[t1, t2] = j2
			if !j2.hasSwitched {
				t.Errorf("should have switched")
			}
		}

		if m1 == m2 {
			t.Errorf("expected different address, got m1: %v, m2: %v\n", m1, m2)
		}

		if j2, ok := m2.(Just[t1, t2]); !ok {
			t.Errorf("expected ok type assert, instead got: %v\n", ok)
		} else {
			if !j2.hasSwitched {
				t.Errorf("type should have switched from customType1 to customType2\n")
			}

			if j2.Next == nil {
				t.Fatalf("j2.Next should not equal nil")
			}
			wantLen := len(testName1)
			got := j2.Next.NameLength
			if wantLen != got {
				t.Errorf("expected NameLength to equal %v, but got %v", wantLen, got)
			}
		}
	})

	t.Run("should avoid nil pointer issues for mapping once between different types", func(t *testing.T) {
		type customType1 struct {
			Name string
			Age  int
		}
		type customType2 struct {
			NameLength int
		}

		vt1 := customType1{
			Name: "test name",
			Age:  39,
		}

		m1 := Just[customType1, customType2]{
			Start: &vt1,
		}

		// m2 has a nil pointer in it
		m2 := m1.Map(func(t *customType1) *customType2 {
			var ret *customType2
			return ret
		}).Map(func(t *customType1) *customType2 {
			// we are protected against nil pointer dereference
			return &customType2{
				NameLength: t.Age,
			}
		})

		if n1, ok := m2.(Just[customType1, customType2]); ok {
			t.Fatalf("n1 should be Nothing, instead got %T", n1)
		}

		if n1, ok := m2.(Nothing[customType1, customType2]); !ok {
			t.Fatalf("n1 should be Nothing, instead got %T", n1)
		}
	})

	t.Run("doing more than 1 maps without a FromMaybeToAnother should be Nothing", func(t *testing.T) {
		type customType1 struct {
			Name string
			Age  int
		}

		myCustomStruct1 := customType1{
			Name: "test name",
			Age:  38,
		}
		m1 := Just[customType1, customType1]{
			Start: &myCustomStruct1,
		}

		m2 := m1.Map(func(t *customType1) *customType1 {
			t.Name = "new test name"
			return t
		})

		if m1 == m2 {
			t.Errorf("expected different address, got m1: %v, m2: %v\n", m1, m2)
		}

		if _, ok := m2.(Just[customType1, customType1]); !ok {
			t.Errorf("expected ok type assert, instead got: %v\n", ok)
		}

		// doesn't matter what we did here, since we didn't reset with FromMaybeToAnother
		// we should get nothing
		m3 := m2.Map(func(t *customType1) *customType1 {
			return t
		})

		if n1, ok := m3.(Nothing[customType1, customType1]); !ok {
			t.Errorf("expected nothing, but got %T", n1)
		}

		if j1, ok := m3.(Just[customType1, customType1]); ok {
			t.Errorf("expected nothing, but got %T", j1)
		}
	})

	t.Run("should be able to map between different types twice using FromMaybeToAnother", func(t *testing.T) {
		type customType1 struct {
			Name string
			Age  int
		}
		type customType2 struct {
			NameLength int
			Age        int
		}
		type customType3 struct {
			NameLengthPlusAge int
		}

		vt1 := customType1{
			Name: "test name",
			Age:  39,
		}

		m1 := Just[customType1, customType2]{
			Start: &vt1,
		}

		m2 := m1.Map(func(t *customType1) *customType2 {
			return &customType2{
				NameLength: len(t.Name),
				Age:        t.Age,
			}
		})

		if m1 == m2 {
			t.Errorf("expected different address, got m1: %v, m2: %v\n", m1, m2)
		}

		if j2, ok := m2.(Just[customType1, customType2]); !ok {
			t.Errorf("expected ok type assert, instead got: %v\n", ok)
		} else {
			if !j2.hasSwitched {
				t.Errorf("type should have switched from customType1 to customType2\n")
			}

			if j2.Next == nil {
				t.Fatalf("j2.Next should not equal nil")
			}
			wantLen := len(testName1)
			got := j2.Next.NameLength
			if wantLen != got {
				t.Errorf("expected NameLength to equal %v, but got %v", wantLen, got)
			}
		}

		m3 := FromMaybeToAnother[customType1, customType2, customType3](m2)
		m4 := m3.Map(func(t *customType2) *customType3 {
			return &customType3{
				NameLengthPlusAge: t.NameLength + t.Age,
			}
		})

		if j3, ok := m4.(Just[customType2, customType3]); !ok {
			t.Errorf("expected Just[customType2, customType3] but instead got %T", j3)
		}

		if n1, ok := m4.(Nothing[customType2, customType3]); ok {
			t.Errorf("expected Just[customType2, customType3] but instead got %T", n1)
		}

	})

	t.Run("should be able to map N times using FromMaybeToAnother", func(t *testing.T) {
		// in progress
		t.SkipNow()
	})
}
