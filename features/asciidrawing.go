package asciiart

func DrawingInput(input []string, bannerSlice []string) (string, error) {
	var convertedStrine string
	for i, str := range input {
		if str == "" && i != 0 {
			convertedStrine += "\n"
			continue
		}
		for h := 1; h < 9; h++ {
			for _, w := range str {
				// if w < 32 || w > 126 {
				// 	log.Fatalf("You Write Non-Printable Char. ")
				// }
				selectChar := int((w-32))*9 + h
				convertedStrine += bannerSlice[selectChar]
			}
			convertedStrine += "\n"
		}
	}
	return convertedStrine, nil
}
