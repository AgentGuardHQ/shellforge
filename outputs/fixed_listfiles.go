func listFiles(params map[string]string, _ int) Result {
	dir := params["directory"]
	if dir == "" {
		dir = params["path"]
	}
	if dir == "" {
		dir = params["dir"]
	}
	if dir == "" {
		dir = "."
	}
	ext := params["extension"]
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if name == "node_modules" || name == ".git" || strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if len(files) > 200 {
			return fmt.Errorf("limit reached")
		}
		if ext != "" && filepath.Ext(name) != ext {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		if rel == "." {
			// Skip the directory itself
			return nil
		}
		if d.IsDir() {
			files = append(files, rel+"/")
		} else {
			files = append(files, rel)
		}
		return nil
	})
	if err != nil {
		return Result{Success: false, Output: err.Error(), Error: "list_error"}
	}
	if len(files) == 0 {
		return Result{Success: true, Output: "(empty directory)"}
	}
	return Result{Success: true, Output: strings.Join(files, "\n")}
}