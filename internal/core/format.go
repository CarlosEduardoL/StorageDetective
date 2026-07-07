package core

// Plural returns word when count is exactly 1, otherwise word with an "s"
// appended. Used for "1 error" vs "2 errors".
func Plural(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

// BoolLabel returns "on" when val is true, "off" when false.
func BoolLabel(val bool) string {
	if val {
		return "on"
	}
	return "off"
}
