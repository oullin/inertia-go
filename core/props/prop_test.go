package props_test

import (
	"testing"

	"github.com/oullin/inertia-go/core/props"
)

func TestAlwaysProp(t *testing.T) {
	p := props.Always("hello")

	if got := p.Prop(); got != "hello" {
		t.Errorf("Always.Prop() = %v, want %q", got, "hello")
	}
}

func TestOptionalProp(t *testing.T) {
	p := props.Optional(42)

	if got := p.Prop(); got != 42 {
		t.Errorf("Optional.Prop() = %v, want 42", got)
	}
}

func TestDeferProp(t *testing.T) {
	t.Run("default group", func(t *testing.T) {
		p := props.Defer("data")

		if p.Group != "default" {
			t.Errorf("Defer.Group = %q, want %q", p.Group, "default")
		}

		if p.IsMerge() {
			t.Error("Defer.IsMerge() = true, want false")
		}

		if got := p.Prop(); got != "data" {
			t.Errorf("Defer.Prop() = %v, want %q", got, "data")
		}
	})

	t.Run("custom group", func(t *testing.T) {
		p := props.Defer("data", "sidebar")

		if p.Group != "sidebar" {
			t.Errorf("Defer.Group = %q, want %q", p.Group, "sidebar")
		}
	})

	t.Run("merge", func(t *testing.T) {
		p := props.Defer("data").Merge()

		if !p.IsMerge() {
			t.Error("Defer.Merge().IsMerge() = false, want true")
		}
	})
}

func TestOnceProp(t *testing.T) {
	p := props.Once("once-val")

	if got := p.Prop(); got != "once-val" {
		t.Errorf("Once.Prop() = %v, want %q", got, "once-val")
	}
}

func TestMergeProp(t *testing.T) {
	t.Run("shallow", func(t *testing.T) {
		p := props.Merge([]int{1, 2})

		if p.IsDeep() {
			t.Error("Merge.IsDeep() = true, want false")
		}
	})

	t.Run("deep", func(t *testing.T) {
		p := props.DeepMerge(map[string]int{"a": 1})

		if !p.IsDeep() {
			t.Error("DeepMerge.IsDeep() = false, want true")
		}
	})
}

func TestScrollProp(t *testing.T) {
	p := props.Scroll([]int{1, 2}, "feedPage", 1, nil, 2).Merge().Reset()

	if !p.IsMerge() {
		t.Error("Scroll.Merge().IsMerge() = false, want true")
	}

	if !p.IsReset() {
		t.Error("Scroll.Reset().IsReset() = false, want true")
	}

	if p.PageName != "feedPage" {
		t.Errorf("Scroll.PageName = %q, want %q", p.PageName, "feedPage")
	}
}
