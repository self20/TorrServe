package torrent

import "server/search/parser"

func GetParser(parserName string) parser.Parser {
	switch parserName {
	case "yohoho":
		return parser.NewYHH()
	case "rutor":
		return parser.NewRutor()
	case "tparser":
		return parser.NewTParser()
	default:
		return nil
	}
}
