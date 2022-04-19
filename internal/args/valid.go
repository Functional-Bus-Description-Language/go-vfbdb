package args

func isValidTarget(target string) bool {
	validTargets := map[string]bool{
		"c-sync": true,
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
	case "c-sync":
		return isValidFlagCSync(flag)
	case "python":
		return isValidFlagPython(flag)
	case "vhdl":
		return isValidFlagVHDL(flag)
	default:
		panic("should never happen")
	}
}

func isValidFlagCSync(flag string) bool {
	validFlags := map[string]bool{
		"-help":       true,
		"-no-asserts": true,
	}

	if _, ok := validFlags[flag]; ok {
		return true
	}

	return false
}

func isValidFlagPython(flag string) bool {
	validFlags := map[string]bool{
		"-help":       true,
		"-no-asserts": true,
	}

	if _, ok := validFlags[flag]; ok {
		return true
	}

	return false
}

func isValidFlagVHDL(flag string) bool {
	validFlags := map[string]bool{
		"-help":   true,
		"-no-psl": true,
	}

	if _, ok := validFlags[flag]; ok {
		return true
	}

	return false
}

func isValidParam(param string, target string) bool {
	if !isValidTarget(target) {
		panic("should never happen")
	}

	switch target {
	case "c-sync":
		return isValidParamCSync(param)
	case "python":
		return isValidParamPython(param)
	case "vhdl":
		return isValidParamVHDL(param)
	}

	return false
}

func isValidParamCSync(param string) bool {
	validParams := map[string]bool{
		"-path": true,
	}

	if _, ok := validParams[param]; ok {
		return true
	}

	return false
}

func isValidParamPython(param string) bool {
	validParams := map[string]bool{
		"-path": true,
	}

	if _, ok := validParams[param]; ok {
		return true
	}

	return false
}

func isValidParamVHDL(param string) bool {
	validParams := map[string]bool{
		"-path": true,
	}

	if _, ok := validParams[param]; ok {
		return true
	}

	return false
}
