package maybe

import (
	"fmt"
	"testing"
)

func TestMaybe_Map(t *testing.T) {
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
			start: &myCustomStruct1,
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
			start: &myCustomStruct1,
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
			start: &vt1,
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

			if j2.next == nil {
				t.Fatalf("j2.Next should not equal nil")
			}
			wantLen := len(testName1)
			got := j2.next.NameLength
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
			start: &vt1,
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

	t.Run("doing more than 1 maps withOUT a FromMaybeToAnother should be Nothing", func(t *testing.T) {
		type customType1 struct {
			Name string
			Age  int
		}

		myCustomStruct1 := customType1{
			Name: "test name",
			Age:  38,
		}
		m1 := Just[customType1, customType1]{
			start: &myCustomStruct1,
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
			start: &vt1,
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

			if j2.next == nil {
				t.Fatalf("j2.Next should not equal nil")
			}
			wantLen := len(testName1)
			got := j2.next.NameLength
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

	t.Run("should be able to map N times using FromMaybeToAnother with valid pointer", func(t *testing.T) {
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
		type customType4 struct {
			NameLengthPlusAgeSq int
		}

		vt1 := customType1{
			Name: "test name",
			Age:  39,
		}

		m1 := Just[customType1, customType2]{
			start: &vt1,
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

			if j2.next == nil {
				t.Fatalf("j2.Next should not equal nil")
			}
			wantLen := len(testName1)
			got := j2.next.NameLength
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

		m5 := FromMaybeToAnother[customType2, customType3, customType4](m4).
			Map(func(t *customType3) *customType4 {
				toRet := &customType4{
					NameLengthPlusAgeSq: t.NameLengthPlusAge * t.NameLengthPlusAge,
				}
				return toRet
			})

		if j4, ok := m5.(Just[customType3, customType4]); !ok {
			t.Errorf("expected Just[customType3, customType4] but instead got %T", j4)
		}

		if n2, ok := m5.(Nothing[customType3, customType4]); ok {
			t.Errorf("expected Just[customType3, customType4] but instead got %T", n2)
		}
	})

	t.Run("should be able to map N times using FromMaybeToAnother with nil pointer", func(t *testing.T) {
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
		type customType4 struct {
			NameLengthPlusAgeSq int
		}

		var vt1 *customType1

		m1 := Of[customType1, customType2](vt1)

		m2 := m1.Map(func(t *customType1) *customType2 {
			// should be protected against nil pointer dereference as this func
			// will never be called
			return &customType2{
				NameLength: len(t.Name),
				Age:        t.Age,
			}
		})

		if j2, ok := m2.(Nothing[customType1, customType2]); !ok {
			t.Errorf("expected Nothing[customType1, customType2], instead got: %T\n", j2)
		}

		m3 := FromMaybeToAnother[customType1, customType2, customType3](m2)
		m4 := m3.Map(func(t *customType2) *customType3 {
			// should be protected against nil pointer dereference as this func
			// will never be called
			return &customType3{
				NameLengthPlusAge: t.NameLength + t.Age,
			}
		})

		if n1, ok := m4.(Nothing[customType2, customType3]); !ok {
			t.Errorf("expected Nothing[customType2, customType3] but instead got %T", n1)
		}

		m5 := FromMaybeToAnother[customType2, customType3, customType4](m4).
			Map(func(t *customType3) *customType4 {
				// should be protected against nil pointer dereference as this func
				// will never be called
				toRet := &customType4{
					NameLengthPlusAgeSq: t.NameLengthPlusAge * t.NameLengthPlusAge,
				}
				return toRet
			})

		if n2, ok := m5.(Nothing[customType3, customType4]); !ok {
			t.Errorf("expected Nothing[customType3, customType4] but instead got %T", n2)
		}
	})

	t.Run("show succinct use of Of, FromMaybeToAnother, and Map with nil pointer", func(t *testing.T) {
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
		type customType4 struct {
			NameLengthPlusAgeSq int
		}

		var none1 *customType1
		otherMaybe1 := Of[customType1, customType2](none1).
			Map(func(c1 *customType1) *customType2 {
				return &customType2{
					NameLength: 0,
					Age:        0,
				}
			})

		otherMaybe2 := FromMaybeToAnother[customType1, customType2, customType3](otherMaybe1).
			Map(func(c2 *customType2) *customType3 {
				return &customType3{
					NameLengthPlusAge: c2.NameLength + c2.Age,
				}
			})

		otherMaybe3 := FromMaybeToAnother[customType2, customType3, customType4](otherMaybe2).
			Map(func(c3 *customType3) *customType4 {
				return nil
			})

		if on1, ok := otherMaybe1.(Nothing[customType1, customType2]); !ok {
			t.Errorf("expected Nothing[customType1, customType2], but got %T", on1)
		}

		if on2, ok := otherMaybe2.(Nothing[customType2, customType3]); !ok {
			t.Errorf("expected Nothing[customType1, customType2], but got %T", on2)
		}

		if on3, ok := otherMaybe3.(Nothing[customType3, customType4]); !ok {
			t.Errorf("expected Nothing[customType1, customType4], but got %T", on3)
		}

		if j1, ok := otherMaybe1.(Just[customType1, customType2]); ok {
			t.Errorf("expected Nothing[customType1, customType2], but got %T", j1)
		}

		if j2, ok := otherMaybe2.(Just[customType2, customType3]); ok {
			t.Errorf("expected Nothing[customType1, customType2], but got %T", j2)
		}

		if j3, ok := otherMaybe3.(Just[customType3, customType4]); ok {
			t.Errorf("expected Nothing[customType1, customType4], but got %T", j3)
		}
	})
}

func TestOf(t *testing.T) {
	t.Run("nil pointer should return Nothing", func(t *testing.T) {
		var n *int
		m1 := Of[int, int](n)

		if n1, ok := m1.(Nothing[int, int]); !ok {
			t.Errorf("expected Nothing but got %T", n1)
		}

		if j1, ok := m1.(Just[int, int]); ok {
			t.Errorf("expected Nothing but got %T", j1)
		}
	})

	t.Run("valid pointer should return Just", func(t *testing.T) {
		n := 42
		m1 := Of[int, int](&n)

		if n1, ok := m1.(Nothing[int, int]); ok {
			t.Errorf("expected Nothing but got %T", n1)
		}

		if j1, ok := m1.(Just[int, int]); !ok {
			t.Errorf("expected Nothing but got %T", j1)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("if nil pointer should return nothing", func(t *testing.T) {
		var n *int
		m1 := Of[int, int](n)
		value := m1.Get()

		if n1, ok := value.(Nothing[int, int]); !ok {
			t.Errorf("expected Nothing but got %T", n1)
		}

		if j1, ok := value.(Just[int, int]); ok {
			t.Errorf("expected Nothing but got %T", j1)
		}
	})

	t.Run("if valid pointer", func(t *testing.T) {
		t.Run("and not yet mapped", func(t *testing.T) {
			const num = 42
			n := num
			m1 := Of[int, int](&n)
			v1 := m1.Get()

			if i1, ok := v1.(*int); !ok {
				t.Errorf("expected *int of value %v, but got type %T", num, i1)
			}

			if i2, ok := v1.(*int); ok {
				if *i2 != num {
					t.Errorf("expected value of %v but got %v", num, *i2)
				}
			}
		})

		t.Run("and HAS been mapped", func(t *testing.T) {
			const num = 42
			n := num
			m1 := Of[int, int](&n).Map(func(i *int) *int {
				value := *i
				toRet := value * value
				return &toRet
			})
			v1 := m1.Get()

			if i1, ok := v1.(*int); !ok {
				t.Errorf("expected *int of value %v, but got type %T", num, i1)
			}

			if i2, ok := v1.(*int); ok {
				if *i2 != num*num {
					t.Errorf("expected value of %v but got %v", num, *i2)
				}
			}
		})
	})
}

func TestAs(t *testing.T) {
	t.Run("should provide sugar over type assertion before mapping", func(t *testing.T) {
		const num = 42
		n := num
		m1 := Of[int, int](&n)
		v1 := m1.Get()

		value, err := As[int](v1)
		if err != nil {
			t.Fatalf("expected nil err but got %v", err)
		}

		if value == nil {
			t.Fatalf("expected NON nil pointer but got nil pointer")
		}

		valueType := fmt.Sprintf("%T", value)
		wantType := "*int"
		if valueType != wantType {
			t.Fatalf("expected %s but instead got %T", wantType, value)
		}

		if *value != num {
			t.Errorf("expected %v but got %v", num, *value)
		}
	})

	t.Run("should provide sugar over type assertion after mapping", func(t *testing.T) {
		const num = 42
		n := num
		want := n + n
		m1 := Of[int, int](&n).Map(func(i *int) *int {
			toRet := *i + *i
			return &toRet
		})
		v1 := m1.Get()

		value, err := As[int](v1)
		if err != nil {
			t.Fatalf("expected nil err but got %v", err)
		}

		if value == nil {
			t.Fatalf("expected NON nil pointer but got nil pointer")
		}

		valueType := fmt.Sprintf("%T", value)
		wantType := "*int"
		if valueType != wantType {
			t.Fatalf("expected %s but instead got %T", wantType, value)
		}

		if *value != want {
			t.Errorf("expected %v but got %v", want, *value)
		}
	})
}
