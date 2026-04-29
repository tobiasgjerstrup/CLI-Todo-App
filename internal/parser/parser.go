package InputParser

import (
	"fmt"
	"strconv"
	"strings"
)

func Parse(clause *string) (string, error) {
	parts := strings.Fields(*clause)

	if len(parts)%4 != 3 {
		return "", fmt.Errorf("where clause is invalid length")
	}

	str, err := stringsToQuery("AND", parts[0], parts[1], parts[2])
	fmt.Println(str)
	if err != nil {
		return "", err
	}
	if len(parts) > 3 {
		for i := 0; i < len(parts); i += 4 {
			if i == 0 {
				i += 3
			}
			fmt.Println(i)
			_str, err := stringsToQuery(parts[i], parts[i+1], parts[i+2], parts[i+3])
			if err != nil {
				return "", err
			}
			str += _str
		}
	}

	return str, nil
}

func stringsToQuery(a string, b string, c string, d string) (string, error) {
	a = strings.ToUpper(a)
	b = strings.ToUpper(b)
	c = strings.ToUpper(c)
	if a != "AND" && a != "OR" {
		return "", fmt.Errorf("expected string to be 'AND' or 'OR', recieved %s", a)
	}

	if b != "ID" && b != "NAME" && b != "STATE" && b != "DESCRIPTION" {
		return "", fmt.Errorf("expected string to be 'ID', 'NAME', 'STATE' or 'DESCRIPTION', recieved %s", b)
	}

	if c != ">" && c != "<" && c != ">=" && c != "<=" && c != "=" && c != "!=" {
		return "", fmt.Errorf("expected string to be '>', '<', '>=', '<=', '=' or '!=', recieved, %s", c)
	}

	if !isNumber(d) {
		return "", fmt.Errorf("expected string to be a valid number, recieved %s", d)
	}

	return fmt.Sprintf("%s %s %s %s ", a, b, c, d), nil
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
