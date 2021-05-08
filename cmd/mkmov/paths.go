package main

import "path"

const (
	OutBaseDir = "out"
)

func videoBaseFilename() string {
	return ActiveProject + "." + string(OutputFileType)
}

func videoOutFilename() string {
	filename := videoBaseFilename()
	return path.Join(OutBaseDir, filename)
}

func projOutDirectory() string {
	// Place to put various outputs for project
	return path.Join("out", ActiveProject+".d")
}

func imageDirectory() string {
	// Directory within project output directory
	outDir := projOutDirectory()
	return path.Join(outDir, "imgs")
}
