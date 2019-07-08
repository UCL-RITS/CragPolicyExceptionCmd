package main

import (
	"log"
	"os/user"
	"strconv"
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
	if currentUser.Username == "SYSTEM" {
		return true
	}

	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		log.Fatal("Could not get current uid: ", err)
	}
	if uid == 0 {
		return true
	}
	// Not 100% sure about this, but I think it's fine
	//  Revised downwards from 1024 for local computers
	if uid < 501 {
		return true
	}

	// Yes I could do some stuff to look this up and convert it
	//  but it's 215.
	// Yes, blame me if this causes a problem.
	//  --Ian
	ccspGid := 215
	gid, err := strconv.Atoi(currentUser.Gid)
	if err != nil {
		log.Fatal("Could not get current gid: ", err)
	}
	if gid == ccspGid {
		return true
	}

	return false
}
