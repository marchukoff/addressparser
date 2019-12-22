package addressparser

import (
	"sort"
)

type parsedAddress struct {
	searchVariant SearchVariant
	lastRecognizedPartIndex int
	notRecognizedPartCount int
	partResults []*ParsedAddressPart
}

func (addr *parsedAddress) addPartResult(part *ParsedAddressPart, partIndex int) {
	addr.partResults = append(addr.partResults, part)
	sort.Slice(addr.partResults, func(i, j int) bool {
		return addr.partResults[i].Level < addr.partResults[j].Level
	})
	if partIndex >= 0 && partIndex > addr.lastRecognizedPartIndex {
		addr.lastRecognizedPartIndex = partIndex
	}
}

func (addr *parsedAddress) getRank() int {
	return len(addr.partResults) - addr.notRecognizedPartCount
}

func (addr *parsedAddress) PostalCode() string {
	if len(addr.partResults) > 0 {
		return addr.partResults[len(addr.partResults)-1].PostalCode
	}
	return ""
}

func (addr *parsedAddress) GetName(level int) string {
	if len(addr.partResults) == 0 {
		return " "
	}

	for _, part := range addr.partResults {
		if part.Level == level{
			return part.Name
		}
	}

	return " "
}

func (addr *parsedAddress) getAddress() *Address {
	address := &Address{
		PostalCode:  addr.PostalCode(),
		Region:      addr.GetName(Region),
		District:    addr.GetName(District),
		Location:    addr.GetName(Location),
		Street:      addr.GetName(Street),
		HouseNumber: addr.GetName(House),
		Building:    "",
		Apartment:   "",
	}

	if address.Location == "" {
		address.Location = addr.GetName(Sublocation)
	}

	apartmentIndex := len(addr.searchVariant.addressParts) - 1
	if addr.lastRecognizedPartIndex < apartmentIndex {
		address.Apartment = addr.searchVariant.addressParts[apartmentIndex]
	}

	address.Building = asString(
		addr.searchVariant.addressParts,
		DefaultHyphens,
		addr.lastRecognizedPartIndex + 1,
		apartmentIndex,
		)

	return address
}
