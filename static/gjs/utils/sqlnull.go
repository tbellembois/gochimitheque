package utils

import "regexp"

// NormalizeSqlNull removes from the obj map keys the .Foo (.String, .Int64...) prefixes
func NormalizeSqlNull(obj map[string]interface{}) *map[string]interface{} {
	// result map
	r := make(map[string]interface{})

	// sqlNull values type detection regexp
	regexS := regexp.MustCompile(`(.+)\.String`)
	regexI := regexp.MustCompile(`(.+)\.Int64`)
	regexF := regexp.MustCompile(`(.+)\.Float64`)
	regexB := regexp.MustCompile(`(.+)\.Bool`)
	regexT := regexp.MustCompile(`(.+)\.Time`)

	for k, iv := range obj {
		// trying to match a regex
		mS := regexS.FindStringSubmatch(k)
		mI := regexI.FindStringSubmatch(k)
		mF := regexF.FindStringSubmatch(k)
		mB := regexB.FindStringSubmatch(k)
		mT := regexT.FindStringSubmatch(k)

		// building the new map without
		// the .Foo in the key names
		if len(mS) > 0 {
			r[mS[1]] = iv.(string)
		} else if len(mI) > 0 {
			r[mI[1]] = iv.(float64)
		} else if len(mF) > 0 {
			r[mF[1]] = iv.(float64)
		} else if len(mB) > 0 {
			r[mB[1]] = iv.(bool)
		} else if len(mT) > 0 {
			r[mT[1]] = iv.(string)
		} else {
			r[k] = iv
		}
	}
	return &r
}
