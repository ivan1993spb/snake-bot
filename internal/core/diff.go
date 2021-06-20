package core

func diff(have, want map[int]int) map[int]int {
	d := make(map[int]int)

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
