package asciiart

func StoreInputAndBanner(Args []string) (string, string) {
	var stringToArt string
	var bannerStyle string
	stringToArt = Args[0]
	bannerStyle = Args[1]
	return stringToArt, bannerStyle
}
