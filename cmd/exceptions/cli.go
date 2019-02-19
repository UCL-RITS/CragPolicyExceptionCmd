package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

// TODO: I feel like there *should* be an edit option, but I'm not sure how it should work.

var (
	// These are set from the CLI in the build command
	commitLabel string
	buildDate   string

	app = kingpin.New("exceptions", "A tool for handling the policy exception entries in the database.")
	// Add debug flag on app here for DB call debugging

	homeDir = os.Getenv("HOME")

	configFile    = app.Flag("config", "Path to config file").Default(homeDir + "/.exceptions_db.conf").String()
	gormDebugMode = app.Flag("ormdebug", "Enable ORM debugging output").Bool()

	listCmd      = app.Command("list", "List entries")
	submitCmd    = app.Command("submit", "Submit a new exception")
	undecideCmd  = app.Command("undecide", "Mark an existing exception as undecided")
	approveCmd   = app.Command("approve", "Approve an existing exception")
	rejectCmd    = app.Command("reject", "Reject an existing exception")
	implementCmd = app.Command("implemented", "Mark an existing exception as implemented")
	removeCmd    = app.Command("remove", "Mark an existing exception as removed")
	formCmd      = app.Command("form", "Handle the exception form files")
	//	editCmd      = app.Command("edit", "Edit an existing exception")
	commentCmd = app.Command("comment", "Add a comment to an existing exception")
	detailsCmd = app.Command("details", "View all details for an exception").Alias("info").Alias("detail")

	jsonDumpCmd   = app.Command("dumpjson", "Full-structured dump of all exceptions as JSON.")
	jsonImportCmd = app.Command("importjson", "Import an array of exceptions as JSON.")

	createDBCmd    = app.Command("createdb", "Create the exceptions DB")
	destroyDBCmd   = app.Command("destroydb", "Destroy the exceptions DB")
	makeNoodlesCmd = app.Command("makenoodles", "Insert some sample data to the database (for development)").Hidden()

	// Need these to make default dates below
	//  Could also make default nil and generate in-function, but have not done that. No particular reason.
	now                     = time.Now()
	nowPlusYear             = time.Now().AddDate(1, 0, 0)
	dateTodayString         = stringFromDate(&now)
	dateTodayPlusYearString = stringFromDate(&nowPlusYear)

	submitName            = submitCmd.Flag("username", "Username exception applies to.").Required().String()
	submitDate            = submitCmd.Flag("submitted", "Date exception was submitted to us. [today]").Default(dateTodayString).String()
	submitStartDate       = submitCmd.Flag("starts", "Date exception should start. [today]").Default(dateTodayString).String()
	submitEndDate         = submitCmd.Flag("ends", "Date exception should finish. [today plus a year]").Default(dateTodayPlusYearString).String()
	submitService         = submitCmd.Flag("service", "Which service the exception applies to.").Default("myriad").String()
	submitExceptionType   = submitCmd.Flag("type", "What type of exception it is.").Default("quota").String()
	submitExceptionDetail = submitCmd.Flag("detail", "Detail of the exception: quota size, queue length, etc.").Default("{'pool': 'home', 'size':'1TB'}").String()

	listOpts      = []string{"all", "undecided", "approved", "rejected", "needed", "active", "removed", "overdue", "pending", "inconsistent", "todo"}
	listHelp      = fmt.Sprintf("Class of exception to list (%s)", strings.Join(listOpts, ", "))
	listClassEnum = listCmd.Arg("class", listHelp).Default("all").Enum(listOpts...)

	attachSubcmd        = formCmd.Command("attach", "Attach a file to an exception.")
	downloadSubcmd      = formCmd.Command("download", "Download a file by file ID.")
	downloadForExSubcmd = formCmd.Command("download-for", "Download all files for an exception.")
	filelistSubcmd      = formCmd.Command("list", "List attached files for an exception.")

	undecideID  = undecideCmd.Arg("id", "").Required().Uint()
	approveID   = approveCmd.Arg("id", "").Required().Uint()
	rejectID    = rejectCmd.Arg("id", "").Required().Uint()
	removeID    = removeCmd.Arg("id", "").Required().Uint()
	implementID = implementCmd.Arg("id", "").Required().Uint()

	undecideForceFlag  = undecideCmd.Flag("force", "Ignore normal transition checks.").Short('f').Bool()
	approveForceFlag   = approveCmd.Flag("force", "Ignore normal transition checks.").Short('f').Bool()
	rejectForceFlag    = rejectCmd.Flag("force", "Ignore normal transition checks.").Short('f').Bool()
	removeForceFlag    = removeCmd.Flag("force", "Ignore normal transition checks.").Short('f').Bool()
	implementForceFlag = implementCmd.Flag("force", "Ignore normal transition checks.").Short('f').Bool()

	attachID      = attachSubcmd.Arg("id", "").Required().Uint()
	downloadID    = downloadSubcmd.Arg("id", "file ID").Required().Uint()
	downloadForID = downloadForExSubcmd.Arg("id", "exception ID").Required().Uint()
	filelistID    = filelistSubcmd.Arg("id", "").Required().Uint()
	//	editID        = editCmd.Arg("id", "").Required().Uint()
	commentID = commentCmd.Arg("id", "").Required().Uint()
	detailsID = detailsCmd.Arg("id", "").Required().Uint()

	// approveApprover = approveCmd.Arg("approver", "Name of the user approving (or 'CRAG')").Required().String()
	// rejectRejecter  = rejectCmd.Arg("rejecter", "Name of the user rejecting (or 'CRAG')").Required().String()
	// //  ^-- Might change these to default to a config file setting later
	// Changed model so that username is always approver or rejecter -- even if CRAG did the actual approving policy-wise

	commentTextArg = commentCmd.Flag("comment", "Comment text -- if not provided, an editor will open for input").Short('c').Default("").String()

	attachFilename = attachSubcmd.Arg("filename", "").Required().String()

	// The other way to do the list suboptions
	//nlistCmd = app.Command("nlist", "List entries")
	//listOpt  = nlistCmd.Arg("type", "list entries of type").Default("all").Enum("all", "submitted", "approved", "needed", "active", "removed", "overdue", "pending", "inconsistent")
)

func main() {
	kingpin.Version(fmt.Sprintf("exceptions commit %s built on %s", commitLabel, buildDate))
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case listCmd.FullCommand():
		list(*listClassEnum)
	case submitCmd.FullCommand():
		submitWithAllParts(*submitName,
			*submitDate,
			*submitStartDate,
			*submitEndDate,
			*submitService,
			*submitExceptionType,
			*submitExceptionDetail)
	case undecideCmd.FullCommand():
		undecide(*undecideID, *undecideForceFlag)
	case approveCmd.FullCommand():
		approve(*approveID, *approveForceFlag)
	case rejectCmd.FullCommand():
		reject(*rejectID, *rejectForceFlag)
	case implementCmd.FullCommand():
		implement(*implementID, *implementForceFlag)
	case removeCmd.FullCommand():
		remove(*removeID, *removeForceFlag)
	case attachSubcmd.FullCommand():
		attach(*attachID, *attachFilename)
	case downloadSubcmd.FullCommand():
		downloadOneFile(*downloadID)
	case downloadForExSubcmd.FullCommand():
		downloadFilesForException(*downloadForID)
	case filelistSubcmd.FullCommand():
		listFilesForException(*filelistID)
		//	case editCmd.FullCommand():
		//		edit(*editID)
	case commentCmd.FullCommand():
		comment(*commentID)
	case detailsCmd.FullCommand():
		details(*detailsID)
	case createDBCmd.FullCommand():
		createDB()
	case destroyDBCmd.FullCommand():
		destroyDB()
	case makeNoodlesCmd.FullCommand():
		makeNoodles()
	case jsonDumpCmd.FullCommand():
		dumpAllAsJson()
	case jsonImportCmd.FullCommand():
		importAllAsJson()
	default:
		kingpin.FatalUsage("Barely-handled error in command-line parsing")
	}
}
