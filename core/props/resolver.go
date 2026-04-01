package props

import (
	"net/http"
	"strings"

	"github.com/oullin/inertia-go/core/httpx"
)

// Result holds the output of prop resolution: the evaluated props,
// the list of merge-prop keys, and the deferred-prop grouping.
type Result struct {
	Props          map[string]any
	MergeProps     []string
	DeepMergeProps []string
	DeferredProps  map[string][]string
	ScrollProps    map[string]ScrollMeta
	OnceProps      map[string]OnceMeta
}

// ScrollMeta is the response metadata for a scrollable prop.
type ScrollMeta struct {
	PageName     string
	PreviousPage any
	NextPage     any
	CurrentPage  any
	Reset        bool
}

// OnceMeta identifies a prop as a once prop on the client.
type OnceMeta struct {
	Prop      string
	ExpiresAt *int64
}

// Resolve filters and evaluates the merged props map according to the
// Inertia.js partial-reload protocol headers found on r.

// propTraits holds metadata discovered by walking a nested prop
// wrapper chain. Each flag indicates whether the corresponding
// wrapper type was found; the associated fields carry its metadata.
type propTraits struct {
	hasAlways   bool
	hasDefer    bool
	hasOptional bool
	hasOnce     bool
	hasScroll   bool
	hasMerge    bool

	deferGroup string
	deferMerge bool

	scrollMeta  ScrollMeta
	scrollMerge bool

	mergeDeep bool

	onceExpiresAt *int64
}

func Resolve(r *http.Request, component string, merged httpx.Props) (*Result, error) {
	res := &Result{
		Props:         make(map[string]any, len(merged)),
		DeferredProps: make(map[string][]string),
		ScrollProps:   make(map[string]ScrollMeta),
		OnceProps:     make(map[string]OnceMeta),
	}

	partialComponent := r.Header.Get(httpx.HeaderPartialComponent)
	isPartial := httpx.IsInertiaRequest(r) && partialComponent == component

	only := splitHeader(r.Header.Get(httpx.HeaderPartialData))
	except := splitHeader(r.Header.Get(httpx.HeaderPartialExcept))
	exceptOnce := splitHeader(r.Header.Get(httpx.HeaderExceptOnceProps))

	onlySet := toSet(only)
	exceptSet := toSet(except)
	exceptOnceSet := toSet(exceptOnce)

	mergeIntent := r.Header.Get(httpx.HeaderInfiniteScroll)

	for key, val := range merged {
		included, err := res.shouldInclude(key, val, isPartial, onlySet, exceptSet, exceptOnceSet, mergeIntent)

		if err != nil {
			return nil, err
		}

		if !included {
			continue
		}

		resolved, err := resolve(val)

		if err != nil {
			return nil, err
		}

		res.Props[key] = resolved
	}

	return res, nil
}

// walkPropChain iterates through the Proper wrapper chain collecting
// trait flags and metadata from every layer. When the same wrapper
// type appears more than once the outermost instance wins.
func walkPropChain(val any) propTraits {
	var t propTraits
	cur := val

	for {
		switch v := cur.(type) {
		case AlwaysProp:
			t.hasAlways = true
			cur = v.Value

		case DeferProp:
			if !t.hasDefer {
				t.hasDefer = true
				t.deferGroup = v.Group
				t.deferMerge = v.IsMerge()
			}

			cur = v.Value

		case OptionalProp:
			t.hasOptional = true
			cur = v.Value

		case OnceProp:
			if !t.hasOnce {
				t.hasOnce = true
				t.onceExpiresAt = v.GetExpiresAt()
			}

			cur = v.Value

		case ScrollProp:
			if !t.hasScroll {
				t.hasScroll = true
				t.scrollMeta = ScrollMeta{
					PageName:     v.PageName,
					PreviousPage: v.PreviousPage,
					NextPage:     v.NextPage,
					CurrentPage:  v.CurrentPage,
					Reset:        v.IsReset(),
				}
				t.scrollMerge = v.IsMerge()
			}

			cur = v.Value

		case MergeProp:
			if !t.hasMerge {
				t.hasMerge = true
				t.mergeDeep = v.IsDeep()
			}

			cur = v.Value

		default:
			return t
		}
	}
}

// shouldInclude determines whether a prop should be included in the
// response and records metadata (merge keys, deferred groups) as a
// side effect. It returns true when the prop should be resolved and
// added to the output.
func (res *Result) shouldInclude(
	key string,
	val any,
	isPartial bool,
	onlySet, exceptSet, exceptOnceSet map[string]struct{},
	mergeIntent string,
) (bool, error) {
	traits := walkPropChain(val)

	// AlwaysProp bypasses all filters.
	if traits.hasAlways {
		res.recordMetadata(key, traits, mergeIntent)

		return true, nil
	}

	// Partial reload filtering.
	if isPartial {
		// exceptSet takes precedence.
		if len(exceptSet) > 0 {
			if _, excluded := exceptSet[key]; excluded {
				return false, nil
			}
		}

		if len(onlySet) > 0 {
			if _, included := onlySet[key]; !included {
				return false, nil
			}
		}
	}

	// DeferProp: excluded on initial, included on partial reload.
	if traits.hasDefer && !isPartial {
		res.DeferredProps[traits.deferGroup] = append(res.DeferredProps[traits.deferGroup], key)

		if traits.deferMerge {
			res.MergeProps = append(res.MergeProps, key)
		}

		return false, nil
	}

	// OptionalProp: excluded on initial, only on explicit request.
	if traits.hasOptional {
		if !isPartial {
			return false, nil
		}

		if len(onlySet) > 0 {
			if _, requested := onlySet[key]; !requested {
				return false, nil
			}
		} else {
			return false, nil
		}
	}

	// OnceProp: record metadata, skip if in except-once set.
	if traits.hasOnce {
		res.OnceProps[key] = OnceMeta{Prop: key, ExpiresAt: traits.onceExpiresAt}

		if _, skip := exceptOnceSet[key]; skip {
			return false, nil
		}
	}

	// Record all remaining metadata for included props.
	res.recordMetadata(key, traits, mergeIntent)

	return true, nil
}

// recordMetadata records scroll, merge, and deferred-merge metadata
// from the collected traits for an included prop.
func (res *Result) recordMetadata(key string, traits propTraits, mergeIntent string) {
	if traits.hasScroll {
		res.ScrollProps[key] = traits.scrollMeta

		if traits.scrollMerge || mergeIntent == "append" {
			res.MergeProps = append(res.MergeProps, key)
		}
	}

	if traits.hasMerge {
		if traits.mergeDeep {
			res.DeepMergeProps = append(res.DeepMergeProps, key)
		} else {
			res.MergeProps = append(res.MergeProps, key)
		}

		return
	}

	// DeferProp merge on partial reload (prop is included).
	if traits.hasDefer && traits.deferMerge {
		res.MergeProps = append(res.MergeProps, key)
	}
}

// splitHeader splits a comma-separated header value into trimmed,
// non-empty tokens.
func splitHeader(val string) []string {
	if val == "" {
		return nil
	}

	parts := strings.Split(val, ",")
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)

		if p != "" {
			out = append(out, p)
		}
	}

	return out
}

func toSet(items []string) map[string]struct{} {
	if len(items) == 0 {
		return nil
	}

	s := make(map[string]struct{}, len(items))

	for _, item := range items {
		s[item] = struct{}{}
	}

	return s
}
