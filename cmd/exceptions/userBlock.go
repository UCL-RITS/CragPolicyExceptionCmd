package main

import (
	"log"
	"os/user"
)

// We could maintain a list, but ccspapp and ccspap2 are
//  both in the ccsp group
// An alternative approach that would be more work would be
//  to look them up in AD or something
func userIsServiceUser() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal("Could not get the current user's details")
	}
	if currentUser.Uid == 0 {
		return true
	}
	if currentUser.Username == "SYSTEM" {
		return true
	}
	if currentUser.Uid < 1024 {
		return true
	}

	// Yes I could do some stuff to look this up and convert it
	//  but it's 215.
	// Yes, blame me if this causes a problem.
	//  --Ian
	ccspGid := 215
	if currentUser.Gid == ccspGid {
		return true
	}

	return false
}
