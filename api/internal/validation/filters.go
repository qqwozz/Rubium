package validation

import (
	"net/url"
	"regexp"
	"strconv"
)

// === Whitelist для enum-полей ===

var (
	validSubjects = map[string]bool{
		"math": true, "physics": true, "cs": true,
		"informatics": true, "chemistry": true, "biology": true,
	}
	validExams = map[string]bool{
		"ege": true, "oge": true, "vpr": true, "olympiad": true,
	}
	validTaskTypes = map[string]bool{
		"choice": true, "number": true, "string": true,
		"multi": true, "code": true, "text": true,
	}
)

// SafeString разрешает только безопасные символы для PostgREST
func SafeString(s string) (string, bool) {
	// Русские/английские буквы, цифры, пробел, дефис, подчёркивание
	re := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]+$`)
	if !re.MatchString(s) {
		return "", false
	}
	return url.QueryEscape(s), true
}

// SafeSearchString для ilike — разрешает ещё wildcard % и _
func SafeSearchString(s string) (string, bool) {
	re := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_%]+$`)
	if !re.MatchString(s) {
		return "", false
	}
	return url.QueryEscape(s), true
}

// ValidSubject проверяет subject по whitelist
func ValidSubject(s string) (string, bool) {
	if !validSubjects[s] {
		return "", false
	}
	return url.QueryEscape(s), true
}

// ValidExam проверяет exam_type по whitelist
func ValidExam(s string) (string, bool) {
	if !validExams[s] {
		return "", false
	}
	return url.QueryEscape(s), true
}

// ValidTaskType проверяет task_type по whitelist
func ValidTaskType(s string) (string, bool) {
	if !validTaskTypes[s] {
		return "", false
	}
	return url.QueryEscape(s), true
}

// ValidDifficulty — число от 1 до 10
func ValidDifficulty(s string) (string, bool) {
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 10 {
		return "", false
	}
	return s, true
}

// ValidTaskNumber — только цифры
func ValidTaskNumber(s string) (string, bool) {
	re := regexp.MustCompile(`^\d+$`)
	if !re.MatchString(s) {
		return "", false
	}
	return url.QueryEscape(s), true
}

// ValidLimit — число от 1 до 100
func ValidLimit(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 100 {
		return 0, false
	}
	return n, true
}

// ValidSort — только разрешённые значения
func ValidSort(s string) (string, bool) {
	switch s {
	case "rating", "newest", "popular":
		return url.QueryEscape(s), true
	default:
		return "", false
	}
}