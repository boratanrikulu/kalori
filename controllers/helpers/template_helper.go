package helpers

func GetTemplateFiles(controller_file string) []string {
	if controller_file == "" {
		panic("Controller file string can't be empty")
	}

	files := []string{
		controller_file,
		"./views/layouts/default.html",
		"./views/partials/head.html",
		"./views/partials/header.html",
		"./views/partials/footer.html",
	}

	return files
}
