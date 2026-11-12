package asset

import "slices"

// ValidateWithStored runs every validation on the incoming rows, treating the
// assets already in the database as part of the hierarchy, and rejects rows
// whose ID is already stored. Both validating and importing use it.
func ValidateWithStored(incoming, existing []Asset) []ValidationError {
	incomingIDs := assetsByID(incoming)

	// Stored assets take part in validation so that parents can be found in
	// the database. Those sharing an ID with an incoming row are left out and
	// reported by checkAlreadyStored instead.
	all := slices.Clone(incoming)
	for _, a := range existing {
		if _, clash := incomingIDs[a.ID]; !clash {
			all = append(all, a)
		}
	}

	errs := []ValidationError{}
	for _, e := range Validate(all) {
		if e.Line != 0 { // line 0 is a stored asset, not a row of the upload
			errs = append(errs, e)
		}
	}
	return append(errs, withType("Already in database", checkAlreadyStored(incoming, existing))...)
}

// checkAlreadyStored reports every incoming row whose asset_id is already in
// the database.
func checkAlreadyStored(incoming, existing []Asset) []ValidationError {
	stored := assetsByID(existing)

	var errs []ValidationError
	for _, a := range incoming {
		if _, exists := stored[a.ID]; exists {
			errs = append(errs, newError(a, "Already Exists", "already exists in the database."))
		}
	}
	return errs
}
