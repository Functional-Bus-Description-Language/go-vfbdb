package args

func isValidTarget(target string) bool {
	validTargets := map[string]bool{
		"python": true,
		"vhdl":   true,
	}

	if _, ok := validTargets[target]; ok {
		return true
	}

	return false
}

func isValidFlag(flag string, target string) bool {
	switch target {
	case "python":
		return isValidFlagPython(flag)
	case "vhdl":
		return isValidFlagVHDL(flag)
	default:
		panic("should never happen")
	}
}

func isValidFlagPython(flag string) bool {
	validFlags := map[string]bool{
		"-h":           true,
		"--help":       true,
		"--no-asserts": true,
	}

	if _, ok := validFlags[flag]; ok {
		return true
	}

	return false
}

func isValidFlagVHDL(flag string) bool {
	validFlags := map[string]bool{
		"-h":       true,
		"--help":   true,
		"--no-psl": true,
	}

	if _, ok := validFlags[flag]; ok {
		return true
	}

	return false
}

func isValidOption(option string, target string) bool {
	if !isValidTarget(target) {
		panic("should never happen")
	}

	switch target {
	case "python":
		return isValidOptionPython(option)
	case "vhdl":
		return isValidOptionVHDL(option)
	}

	return false
}

func isValidOptionPython(option string) bool {
	validOptions := map[string]bool{
		"--path": true,
	}

	if _, ok := validOptions[option]; ok {
		return true
	}

	return false
}

func isValidOptionVHDL(option string) bool {
	validOptions := map[string]bool{
		"--path": true,
	}

	if _, ok := validOptions[option]; ok {
		return true
	}

	return false
}
