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

	// AlwaysProp bypasses all filters.
	if _, ok := val.(AlwaysProp); ok {
		res.recordMerge(key, val)

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

	switch v := val.(type) {
	case DeferProp:
		if !isPartial {
			// Record for client-side lazy loading.
			res.DeferredProps[v.Group] = append(res.DeferredProps[v.Group], key)

			if v.IsMerge() {
				res.MergeProps = append(res.MergeProps, key)
			}

			return false, nil
		}

		res.recordMerge(key, val)

		return true, nil

	case OptionalProp:
		if !isPartial {
			return false, nil
		}

		// Only include it if explicitly requested.
		if len(onlySet) > 0 {
			_, requested := onlySet[key]

			return requested, nil
		}

		return false, nil

	case OnceProp:
		res.OnceProps[key] = OnceMeta{Prop: key}

		if _, skip := exceptOnceSet[key]; skip {
			return false, nil
		}

		return true, nil

	case ScrollProp:
		res.ScrollProps[key] = ScrollMeta{
			PageName:     v.PageName,
			PreviousPage: v.PreviousPage,
			NextPage:     v.NextPage,
			CurrentPage:  v.CurrentPage,
			Reset:        v.IsReset(),
		}

		if v.IsMerge() || mergeIntent == "append" {
			res.MergeProps = append(res.MergeProps, key)
		}

		return true, nil

	case MergeProp:
		res.recordMerge(key, val)

		return true, nil

	default:
		// Regular prop: on a non-partial request always include.
		return true, nil
	}
}

// recordMerge adds the key to MergeProps or DeepMergeProps when val
// is a merge-type prop.
func (res *Result) recordMerge(key string, val any) {
	switch v := val.(type) {
	case MergeProp:
		if v.IsDeep() {
			res.DeepMergeProps = append(res.DeepMergeProps, key)
		} else {
			res.MergeProps = append(res.MergeProps, key)
		}
	case DeferProp:
		if v.IsMerge() {
			res.MergeProps = append(res.MergeProps, key)
		}
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
