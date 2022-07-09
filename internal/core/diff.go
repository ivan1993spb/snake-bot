package core

// TODO: Add checks for negative numbers.

type State map[uint]uint

type Diff map[int]int

func diff(have, want State) Diff {
	d := make(Diff)

	for k, v1 := range have {
		if v2, ok := want[k]; ok {
			if v := v2 - v1; v != 0 {
				d[k] = v
			}
		} else {
			d[k] = -v1
		}
	}

	for k, v := range want {
		if _, ok := d[k]; ok {
			continue
		}
		if _, ok := have[k]; ok {
			continue
		}
		d[k] = v
	}

	return d
}

func diffOne(k, v int, have State) Diff {
	return Diff{
		k: v - have[k],
	}
}
