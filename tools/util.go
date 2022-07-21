package tools

import "strconv"

func ToInt(t interface{}) int {
	switch t := t.(type) {
	case int:
		return t
	case string:
		tt, _ := strconv.Atoi(t)
		return tt
	case float64:
		return int(t)
	case interface{}:
		return t.(int)
	}

	return 0
}

func ToString(t interface{}) string {
	switch t := t.(type) {
	case int:
		return strconv.Itoa(t)
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', 2, 64)
	case interface{}:
		return t.(string)
	}

	return ""
}

func ToFloat64(t interface{}) float64 {
	switch t := t.(type) {
	case int:
		return float64(t)
	case string:
		tt, _ := strconv.ParseFloat(t, 64)
		return tt
	case float64:
		return t
	case interface{}:
		return t.(float64)
	}

	return 0.00
}

func ToFloat32(t interface{}) float32 {
	switch t := t.(type) {
	case int:
		return float32(t)
	case string:
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		tt, _ := strconv.ParseFloat(t, 64)
=======
		tt, _ := strconv.ParseFloat(t, 32)
>>>>>>> 5fa5f02c1373b226cd4ab46bcdfa3326f6ae89d0
=======
		tt, _ := strconv.ParseFloat(t, 32)
>>>>>>> 5fa5f02c1373b226cd4ab46bcdfa3326f6ae89d0
=======
		tt, _ := strconv.ParseFloat(t, 32)
>>>>>>> 5fa5f02c1373b226cd4ab46bcdfa3326f6ae89d0
		return float32(tt)
	case float32:
		return t
	case interface{}:
		return t.(float32)
	}

	return 0.00
}
