package main

import (
	"fmt"
	"strings"
)

// Gives a summary for the week of current state and things that will need to be done.
func report() {
	// Note: this produces YAML, but, really *dumb* YAML that looks like normal text
	// Otherwise we could just use yaml.Marshal
	reportData := gatherReportData()
	rs := "Report:\n" // Report string
	for catName, longCatName := range map[string]string{
		"decision waiting":         "Waiting for Decision",
		"implementation waiting":   "Waiting for Implementation",
		"removal waiting":          "Waiting for Removal",
		"expires within five days": "Will Expire Within Five Days",
		"expires within two weeks": "Will Expire Within Two Weeks",
	} {

		ids := reportData[catName]
		c := len(ids)
		if c != 0 {
			rs = rs + fmt.Sprintf("  %s:\n    Count: %d\n    IDs:\n%s", longCatName, c, yamlListOfIDs(ids, 5))
		}
	}
	fmt.Print(rs)
}

func yamlListOfIDs(list []uint, indent int) string {
	var s string
	sIndent := strings.Repeat(" ", indent)
	for _, v := range list {
		s = s + fmt.Sprintf("%s- %d\n", sIndent, v)
	}
	return s
}

func gatherReportData() map[string][]uint {
	timeNow := getDBTimeNow()
	time5daysFromNow := getDBFutureTimes("5 day")
	timeTwoWeeksFromNow := getDBFutureTimes("14 day")

	rd := make(map[string][]uint)
	// We want the count of and IDs of:
	//   - things that are currently waiting for a decision
	rd["decision waiting"] = getExceptionIDsWhere("status = 'undecided'")
	//   - things that are currently waiting for implementation
	rd["implementation waiting"] = getExceptionIDsWhere("status = 'approved' AND start_date > " + timeNow)
	//   - things that are waiting to be removed
	rd["removal waiting"] = getExceptionIDsWhere("status = 'implemented' AND end_date < " + timeNow)
	//   - things that will expire within the next 5 days (ie. working week)
	rd["expires within five days"] = getExceptionIDsWhere("status = 'implemented' AND end_date >= " + timeNow + " AND end_date < " + time5daysFromNow)
	//   - things that will expire in the next 5-14 days (ie. their owner should be notified)
	rd["expires within two weeks"] = getExceptionIDsWhere("status = 'implemented' AND end_date >= " + time5daysFromNow + " AND end_date < " + timeTwoWeeksFromNow)

	return rd
}

func getExceptionIDsWhere(whereClause string) []uint {
	exs := getExceptionsWhere(whereClause)
	exIDs := make([]uint, 0)
	for _, v := range exs {
		exIDs = append(exIDs, v.ID)
	}
	return exIDs
}

func getExceptionsWhere(whereClause string) []Exception {
	var listSet []Exception
	db := getDB()
	defer db.Close()

	db.Where(whereClause).Find(&listSet)
	return listSet
}
