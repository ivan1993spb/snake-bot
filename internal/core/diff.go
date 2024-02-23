package core

// diff returns the difference between two maps.
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

func stateBotsNumber(state map[int]int) int {
	number := 0
	for _, bots := range state {
		number += bots
	}
	return number
}

func diffStats(m map[int]int) (int, int) {
	var add, remove int
	for _, v := range m {
		if v > 0 {
			add += v
		} else {
			remove -= v
		}
	}
	return add, remove
}
