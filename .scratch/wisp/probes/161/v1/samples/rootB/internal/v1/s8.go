package v1

func pull(mirrorURL, sha256Sum string) string {
	return mirrorURL + sha256Sum
}
