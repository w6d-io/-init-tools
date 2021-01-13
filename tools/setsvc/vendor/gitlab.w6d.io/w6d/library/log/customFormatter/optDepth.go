package runtime

import (
	"strings"
)

// GetDepthFunction ...
func (f *Formatter) GetDepthFunction(pathFunc string, pathFile string) (string, string) {

	depth := f.FuncDepth
	indexLastFolder := strings.LastIndex(pathFunc, "/") + 1
	if f.PathFile {
		index := strings.LastIndex(pathFunc, "/")
		if index != -1 {
			pathFile = pathFunc[:index+strings.Index(pathFunc[index:], ".")] + "/" + pathFile
		} else {
			pathFile = pathFunc + "/" + pathFile
		}
	}

	if strings.Contains(pathFunc[indexLastFolder:], "func") {
		depth += strings.Count(pathFunc[indexLastFolder:], ".")
	}

	if strings.Contains(pathFunc, ".init.0") {
		return "init", pathFile
	}

	maxDepth := strings.Count(pathFunc[indexLastFolder:], ".")
	pathFunc = pathFunc[indexLastFolder:]
	if depth >= maxDepth {
		maxDepth = 1
	} else {
		maxDepth = maxDepth - depth
	}
	for i := 0; i < maxDepth; i++ {
		pathFunc = pathFunc[strings.Index(pathFunc, ".")+1:]
	}
	return pathFunc, pathFile
}
