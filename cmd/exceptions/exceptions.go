package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/olekukonko/tablewriter"
)

type Exception struct {
	// Book-keeping and bureaucracy tracking
	gorm.Model
	Username        string     `gorm:"type:varchar(10);not null"`
	SubmittedDate   *time.Time `gorm:"default:NULL"`
	StartDate       *time.Time `gorm:"default:NULL"`
	EndDate         *time.Time `gorm:"default:NULL"`
	Service         string     `gorm:"type:varchar(16);not null"`
	ExceptionType   string     `gorm:"type:varchar(128);not null"`
	ExceptionDetail string     `gorm:"type:varchar(512);not null"`
	FormFiles       []FormFile
	Comments        []Comment
	StatusChanges   []StatusChange
}

type FormFile struct {
	gorm.Model
	Exception    Exception
	ExceptionID  uint
	FileName     string
	FileContents []byte `gorm:"mediumblob"`
}

type Comment struct {
	gorm.Model
	Exception   Exception
	ExceptionID uint
	CommentBy   string `gorm:"type:varchar(10); not null"`
	CommentText string
}

type StatusChange struct {
	gorm.Model
	Exception   Exception
	ExceptionID uint
	OldStatus   string `gorm:"type:varchar(16);default:'none';not null"`
	NewStatus   string `gorm:"type:varchar(16);default:'none';not null"`
	Changer     string `gorm:"type:varchar(10); not null"`
}

// This is our half-assed state machine map doodad. It determines
//  what states can be made.
func isValidChange(oldStatus string, newStatus string) bool {
	validChanges := make(map[string][]string)
	validChanges["(none)"] = []string{"undecided"}
	validChanges["undecided"] = []string{"approved", "rejected"}
	validChanges["approved"] = []string{"implemented"}
	validChanges["implemented"] = []string{"removed"}
	validChanges["removed"] = []string{}
	validChanges["rejected"] = []string{}

	for _, v := range validChanges[oldStatus] {
		if v == newStatus {
			return true
		}
	}
	return false
}

func GetException(id uint) *Exception {
	db := getDB()
	exception := &Exception{}
	db.First(&exception, id)
	return exception
}

func (exception *Exception) GetStatus() string {
	db := getDB()
	if db.Model(exception).Association("StatusChanges").Count() == 0 {
		return "(none)"
	}

	var statusChanges []StatusChange

	// Iterates over all the status changes for an entry and finds
	//  the one with the largest created time, to get the most recent
	db.Model(exception).Related(&statusChanges)
	maxTime := time.Unix(0, 0)
	maxK := 0
	for k, v := range statusChanges {
		if maxTime.Before(v.CreatedAt) {
			maxK = k
			maxTime = v.CreatedAt
		}
	}
	return statusChanges[maxK].NewStatus
}

func (exception *Exception) ChangeStatusTo(newStatus string) error {
	currentStatus := exception.GetStatus()
	if !isValidChange(currentStatus, newStatus) {
		return errors.New(fmt.Sprintf("Proposed status change (%s -> %s) is invalid", currentStatus, newStatus))
	}

	currentUser, err := user.Current()
	if err != nil {
		panic("Could not get the current username")
	}
	statusChange := &StatusChange{ExceptionID: exception.ID, OldStatus: currentStatus, NewStatus: newStatus, Changer: currentUser.Username}

	db := getDB()
	db.Create(statusChange)
	db.NewRecord(statusChange)
	db.Save(exception)

	return nil
}

func (exception *Exception) AddComment(text string) {

	currentUser, err := user.Current()
	if err != nil {
		panic("could not get current user name")
	}
	currentUsername := currentUser.Username
	comment := &Comment{ExceptionID: exception.ID, CommentText: text, CommentBy: currentUsername}

	db := getDB()
	db.Save(comment)
}

func (exception *Exception) DurationRemaining() (*time.Duration, string) {
	if exception.EndDate == nil || exception.StartDate == nil {
		return nil, "--"
	}

	if time.Now().After(*exception.EndDate) {
		return nil, "finished"
	}

	if time.Now().Before(*exception.StartDate) {
		return nil, "not started yet"
	}

	duration := time.Until(*exception.EndDate)
	return &duration, ""
}

func approve(ID uint) {
	exception := GetException(ID)
	err := exception.ChangeStatusTo("approved")
	if err != nil {
		fmt.Println(err)
	}
}

func reject(ID uint) {
	exception := GetException(ID)
	err := exception.ChangeStatusTo("rejected")
	if err != nil {
		fmt.Println(err)
	}
}

func implement(ID uint) {
	exception := GetException(ID)
	err := exception.ChangeStatusTo("implemented")
	if err != nil {
		fmt.Println(err)
	}
}

func remove(ID uint) {
	exception := GetException(ID)
	err := exception.ChangeStatusTo("removed")
	if err != nil {
		fmt.Println(err)
	}
}

func notYetImplemented() {
	fmt.Printf("This thing not yet implemented.\n")
	panic("!")
}

func list(kind string) {
	//timeNow := "NOW()" // MySQL
	timeNow := "date('now')" // SQLite

	// In theory you'd use these to determine unset but you can just use zero instead
	//zeroTime := "FROM_UNIXTIME(0)" // MySQL
	//zeroTime := "date(0, 'unixepoch')" // SQLite (I think)

	db := getDB()
	defer db.Close()
	var listSet []Exception
	switch kind {
	case "all":
		db.Find(&listSet)
		printExceptionTableSummary(listSet)
	case "pending":
		db.Where(timeNow + " < start_date AND decided_date IS NOT NULL").Find(&listSet)
		printExceptionTableSummary(listSet)
	case "undecided":
		db.Where("status = 'undecided'").Find(&listSet)
		printExceptionTableSummary(listSet)
	case "approved":
		db.Where("status = 'approved'").Find(&listSet)
		printExceptionTableSummary(listSet)
	case "needed":
		db.Where("status = 'approved' AND start_date > " + timeNow).Find(&listSet)
		printExceptionTableSummary(listSet)
	case "active":
		db.Where("status = 'implemented'").Find(&listSet)
		printExceptionTableSummary(listSet)
	case "overdue":
		db.Where("status = 'implemented' AND end_date < " + timeNow).Find(&listSet)
		printExceptionTableSummary(listSet)
	case "removed":
		db.Where("status = 'removed'").Find(&listSet)
		printExceptionTableSummary(listSet)
	case "todo":
		db.Where("(status = 'implemented' AND end_date < " + timeNow + ") OR (status = 'approved' AND start_date > " + timeNow + ") OR (status = 'undecided')").Find(&listSet)
		printExceptionTableSummary(listSet)
	case "inconsistent":
		// This should cover any weird states that an exception could get into where it needs fixing
		// TODO try to think of more states here
		// TODO move this out into a call like IsInconsistent and then run for each Exception
		db.Where("(submitted_date IS NULL) OR " +
			"(start_date IS NULL AND end_date IS NOT NULL) OR " +
			"(removed_date IS NOT NULL AND implemented_date IS NULL) OR " +
			"(implemented_date IS NOT NULL AND submitted_date IS NULL) OR " +
			"(removed_date < implemented_date) OR (removed_date < decided_date) OR (removed_date < submitted_date) OR " +
			"(implemented_date < decided_date) OR (implemented_date < submitted_date) OR " +
			"(decided_date < submitted_date) OR" +
			"(start_date > end_date) OR " +
			"((state = 'undecided') AND ((decided_date IS NOT NULL) OR (implemented_date IS NOT NULL) OR (removed_date IS NOT NULL))) OR " +
			"((state = 'approved') AND ((implemented_date IS NOT NULL) OR (removed_date IS NOT NULL))) OR " +
			"((state = 'implemented') AND (removed_date IS NOT NULL)) " +
			"").Find(&listSet)
		printExceptionTableSummary(listSet)
	default:
		notYetImplemented()
	}
}

//var epochZero = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

func stringFromDate(timeIn *time.Time) string {
	if timeIn == nil {
		return "--"
	} else {
		return timeIn.Format("2006-01-02")
	}
}

func printExceptionTableSummary(exceptions []Exception) {
	db := getDB()

	if len(exceptions) == 0 {
		fmt.Println("No such records found.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Username", "Status", "Sub Date", "Start Date", "End Date", "Type", "Detail", "Attachments", "Comments"})
	table.SetBorder(false)

	for _, ex := range exceptions {
		var statusString string
		numComments := db.Model(&ex).Association("Comments").Count()
		numAttachments := db.Model(&ex).Association("FormFiles").Count()
		statusString = ex.GetStatus()
		table.Append([]string{fmt.Sprintf("%d", ex.ID),
			ex.Username,
			statusString,
			stringFromDate(ex.SubmittedDate),
			stringFromDate(ex.StartDate),
			stringFromDate(ex.EndDate),
			ex.ExceptionType,
			ex.ExceptionDetail,
			fmt.Sprintf("%d", numAttachments),
			fmt.Sprintf("%d", numComments),
		})
	}
	table.Render()
}

func submitWithAllParts(username string, submitDateString string, startDateString string, endDateString string, service string, exceptionType string, details string) {
	// First convert dates into proper formats
	var submitDate time.Time
	var startDate time.Time
	var endDate time.Time
	var err error

	submitDate, err = time.Parse("2006-01-02", submitDateString)
	if err != nil {
		panic("Could not parse submit date")
	}

	startDate, err = time.Parse("2006-01-02", startDateString)
	if err != nil {
		panic("Could not parse start date")
	}

	endDate, err = time.Parse("2006-01-02", endDateString)
	if err != nil {
		panic("Could not parse end date")
	}

	// Then create the exception
	exception := Exception{Username: username,
		SubmittedDate:   &submitDate,
		StartDate:       &startDate,
		EndDate:         &endDate,
		Service:         service,
		ExceptionType:   exceptionType,
		ExceptionDetail: details}

	db := getDB()
	defer db.Close()
	db.NewRecord(exception)
	db.Create(&exception)

	exception.ChangeStatusTo("undecided")
}

func comment(id uint) {
	db := getDB()
	exception := &Exception{}

	exRetrErrors := db.First(&exception, id).GetErrors()

	if exception.ID == 0 {
		fmt.Println("No record of that exception.")
		return
	}

	for _, v := range exRetrErrors {
		fmt.Println(v)
	}

	// commentTextArg is a command line option set in cli.go
	var commentText string
	var err error
	if *commentTextArg == "" {
		commentText, err = getTextFromEditor()
		if err != nil {
			panic(err)
		}
	} else {
		commentText = *commentTextArg
	}

	currentUser, err := user.Current()
	currentUsername := currentUser.Username
	comment := &Comment{ExceptionID: id, CommentText: commentText, CommentBy: currentUsername}

	db.Save(comment)
	return
}

func timeRemaining(exception *Exception) string {
	remaining, msg := exception.DurationRemaining()

	if remaining == nil {
		return msg
	}

	return fmt.Sprintf("%d days", int(remaining.Hours()/24))
}

func details(id uint) {
	db := getDB()
	exception := &Exception{}
	var comments []Comment
	var files []FormFile
	var statusChanges []StatusChange

	errors := db.First(&exception, id).GetErrors()

	if exception.ID == 0 {
		fmt.Println("No record of that exception.")
		return
	}

	for _, v := range errors {
		fmt.Println(v)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(80)

	timeRemaining := timeRemaining(exception)

	data := [][]string{
		[]string{"ID", fmt.Sprint(exception.ID)},
		[]string{"Username", exception.Username},
		[]string{"Type", exception.ExceptionType},
		[]string{"Detail", exception.ExceptionDetail},
		[]string{"Created", exception.CreatedAt.Format("2006-01-02 15:04:05 MST")},
		[]string{"Updated", exception.UpdatedAt.Format("2006-01-02 15:04:05 MST")},
		[]string{"Submitted", stringFromDate(exception.SubmittedDate)},
		[]string{"Starts", stringFromDate(exception.StartDate)},
		[]string{"Ends", stringFromDate(exception.EndDate)},
		[]string{"Remaining", timeRemaining},
		[]string{"Status", exception.GetStatus()},
	}

	db.Model(&exception).Related(&statusChanges)
	if len(statusChanges) == 0 {
		data = append(data, []string{"Status Updates", "(none)"})
	} else {
		statusRowLabel := "Status Change"
		for _, v := range statusChanges {
			data = append(data, []string{statusRowLabel, fmt.Sprintf("%s -> %s, by %s [%s]", v.OldStatus, v.NewStatus, v.Changer, v.UpdatedAt.Format("2006-01-02"))})
		}
	}

	db.Model(&exception).Related(&files)
	if len(files) == 0 {
		data = append(data, []string{"File", "(none)"})
	} else {
		fileRowLabel := "File"
		for _, v := range files {
			data = append(data, []string{fileRowLabel, fmt.Sprintf("%s (%d bytes)", v.FileName, len(v.FileContents))})
			fileRowLabel = ""
		}
	}

	db.Model(&exception).Related(&comments)
	if len(comments) == 0 {
		data = append(data, []string{"Comment", "(none)"})
	} else {
		commentRowLabel := "Comment"
		for _, v := range comments {
			data = append(data, []string{commentRowLabel, v.CommentText})
			commentRowLabel = ""
		}
	}

	table.AppendBulk(data)
	table.Render()

}
