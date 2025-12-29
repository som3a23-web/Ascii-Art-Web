package asciiart

func StoreInputAndBanner(Args []string) (string, string) {
	var stringToArt string
	var bannerStyle = "standard"

	switch len(Args) {
	case 3:
		stringToArt = Args[1]
		bannerStyle = Args[2]
	case 2:
		stringToArt = Args[1]
	default:
		stringToArt = ""
	}

	return stringToArt, bannerStyle
}
