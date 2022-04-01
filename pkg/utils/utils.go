package utils

import (
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"
)

func SnakeCasePath(path string) string {
	re := regexp.MustCompile(`[[:punct:]](\w)`)
	p := re.ReplaceAllStringFunc(path, strings.Title)
	re = regexp.MustCompile(`[[:punct:]]`)
	return re.ReplaceAllString(p, "")
}
func CreateModFormatted(modName, version, goVersion, replace string) ([]byte, error) {
	mod := modfile.File{}
	err := mod.AddModuleStmt(modName)
	if err != nil {
		return []byte{}, err
	}

	err = mod.AddGoStmt(goVersion)
	if err != nil {
		return []byte{}, err
	}

	err = mod.AddRequire("github.com/alxbse/toiletpaper", version)
	if err != nil {
		return []byte{}, err
	}

	if replace != "" {
		err = mod.AddReplace("github.com/alxbse/toiletpaper", version, replace, "")
		if err != nil {
			return []byte{}, err
		}
	}

	modFormatted, err := mod.Format()
	if err != nil {
		return []byte{}, err
	}

	return modFormatted, nil
}
