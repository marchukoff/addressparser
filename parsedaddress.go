package addressparser

type parsedAddress struct {
	searchVariant SearchVariant
	lastRecognizedPartIndex int
	notRecognizedPartCount int
	partResults []*ParsedAddressPart
}

