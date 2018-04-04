package utils

import "regexp"

func FileToLink(file string) string {
	re := regexp.MustCompile(`[ !\*'\(\);:@&=\+\$,/\?#\[\]~",]`)
	return re.ReplaceAllString(file, `_`)
}
