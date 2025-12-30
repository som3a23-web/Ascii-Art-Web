package asciiart

func DrawingInput(input []string, bannerSlice []string) string {
	var convertedStrine string
	for _, str := range input {
		if str == "" {
			convertedStrine += "\n"
			continue
		}
		for h := 1; h < 9; h++ {
			for _, w := range str {
				selectChar := int((w-32))*9 + h
				convertedStrine += bannerSlice[selectChar]

			}
			convertedStrine += "\n"
		}
	}
	return convertedStrine
}
