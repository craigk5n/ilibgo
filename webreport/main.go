//  Report generator for apache2 access.log file
//
//  Usage:
//    webreport [options] [-tod|-dom|-dow|-moy] accesslogfile [accesslogfile ...]
//
//    Default action is to generate a day of the month usage graph to stdout.
//
//    -tod                Display a time of day report indicating what
//                        hours of the day get the most usage.
//                        [THIS IS THE DEFAULT]
//    -dom                Display a day of month report indicating what
//                        days of the month get the most usage.
//    -dow                Display a day of week report indicating what
//                        days of the week get the most usage.
//    -moy                Displays month report indicating usage by month
//
//    options             what it does
//    ----------------    ------------------------------------------------
//    -all                Use all data [default]
//    -today              Use data from today only
//    -yesterday          Use data from yesterday only
//    -lastweek           Use data from the Mon-Sun week prior to this one
//    -thisweek           Use data for this Mon-Sun week
//    -thismonth          Use data for the the current month
//    -lastmonth          Use data for the month prior to this one
//    -bar                Use a bar graph instead of a line graph.
//    -line               Use line graph [default]
//    -nohdr              Do not display title and summary at bottom
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/craigk5n/ilibgo"
	font "github.com/craigk5n/ilibgo/fonts/adobe_utopia_100dpi"
)

// History:
//  18-Aug-2022   Craig Knudsen        craig@k5n.us
//                Converted from C to Go
//  29-May-1996   Craig Knudsen        cknudsen@radix.net
//                Converted from gd library to Ilib.
//  21-Sep-1994   Craig Knudsen        cknudsen@radix.net
//                Created
//
// TODO:
// - Filter logs by HTTP code (ignore 302, 404, etc)
// - Use IP to calculate unique visitors
// - Use IP to create chart of countries, cities etc.  (See https://github.com/cc14514/go-geoip2-db)
// - Better font choices
// - Custom date ranges passed in from command line params
// - Support additional log formats
// - I18N

const (
	LeftPad    int = 80
	TopPad     int = 50
	RightPad   int = 50
	BottomPad  int = 55
	DataWidth  int = 25
	DataHeight int = 300
)

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// Define start and stop times.  Use YYMMDDHHMMSS so that we can just
// use == to compare times.
var startTime time.Time /* in Time with local timezone */
var stopTime time.Time  /* in Time with local timezone */
var prettyStart string  /* in human readable format */
var prettyStop string   /* in human readable format */
var months []string = []string{
	"Jan", "Feb", "Mar", "Apr", "May", "Jun",
	"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
}
var wdays []string = []string{
	"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
}
var daysInMonth []int = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
var daysInLeapMonth []int = []int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
var tod [24]int /* count by time of day */
var dow [7]int  /* couny by day of week */
var dom [31]int /* count by day of month */
var moy [12]int /* count by month */

const Sunday int = 0
const Monday int = 1
const Tuesday int = 2
const Wednesday int = 3
const Thursday int = 4
const Friday int = 5
const Saturday int = 6

var total int = 0 /* total number of requests counted */

type OutputType int

const DayOfMonth OutputType = 1
const DayOfWeek OutputType = 2
const TimeOfDay OutputType = 3
const MonthOfYear OutputType = 4

type GraphType int

const BarGraph GraphType = 1
const LineGraph GraphType = 2

var outputType OutputType = OutputType(DayOfMonth)
var graphType GraphType = GraphType(DayOfMonth)
var displayHeader bool = true
var lineWidth int = 2 // for line graphs (can use 1-3)

// String to use in place of "Views"
var yAxisLabel string = "Views"

func main() {
	// default is to use all data
	startTime = time.Now().AddDate(-100, 0, 0) // 100 years ago
	stopTime = time.Now().AddDate(1, 0, 0)     // next year
	outfile := ""
	title := "Web Report 2.0"

	// process command line arguments
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-bar" {
			graphType = BarGraph
		} else if arg == "-line" {
			graphType = LineGraph
		} else if arg == "-tod" || arg == "-timeofday" {
			outputType = TimeOfDay
		} else if arg == "-dom" || arg == "-dayofmonth" {
			outputType = DayOfMonth
		} else if arg == "-dow" || arg == "-dayofweek" {
			outputType = DayOfWeek
		} else if arg == "-moy" || arg == "-monthofyear" {
			outputType = MonthOfYear
		} else if arg == "-t" || arg == "-today" {
			setTimesToday()
		} else if arg == "-y" || arg == "-yesterday" {
			setTimesYesterday()
		} else if arg == "-lw" || arg == "-lastweek" {
			setTimesLastWeek()
		} else if arg == "-tw" || arg == "-thisweek" {
			setTimesThisWeek()
		} else if arg == "-tm" || arg == "-thismonth" {
			setTimesThisMonth()
		} else if arg == "-lm" || arg == "-lastmonth" {
			setTimesLastMonth()
		} else if arg == "-nohdr" || arg == "-noheader" {
			displayHeader = false
		} else if arg == "-yaxislabel" {
			if i+1 >= len(os.Args) {
				fmt.Printf("Parameter -yaxislabel requres an integer parameter")
				os.Exit(1)
			}
			i++
			yAxisLabel = os.Args[i]
		} else if arg == "-o" || arg == "-outfile" {
			if i+1 >= len(os.Args) {
				fmt.Printf("Parameter -outfile requres an output filename")
				os.Exit(1)
			}
			i++
			outfile = os.Args[i]
		} else if arg == "-t" || arg == "-title" {
			if i+1 >= len(os.Args) {
				fmt.Printf("Parameter -title requres an string")
				os.Exit(1)
			}
			i++
			title = os.Args[i]
		} else if strings.HasPrefix(arg, "-") {
			fmt.Printf("unrecognized parameter")
			os.Exit(1)
		} else {
			readFile(arg)
		}
	}

	// create output
	if len(outfile) == 0 {
		fmt.Printf("You must specify an output file with -outfile\n")
		os.Exit(1)
	}
	fp, err := os.Create(outfile)
	if err != nil {
		fmt.Printf("Error writing to %s: %v\n", outfile, err)
		os.Exit(1)
	}
	generateOutput(fp, title)
	fp.Close()
}

func setTimes(start time.Time, end time.Time) {
	startTime = start
	stopTime = end
}

// Set the time range variables so that we look at data from today only.
func setTimesToday() {
	t := time.Now()
	setTimes(t, t)
}

// Set the time range variables so that we look at data from yesterday only.
func setTimesYesterday() {
	setTimes(time.Now(), time.Now().AddDate(0, 0, -1))
}

// Set the time range variables so that we look at data from last week
// where the week starts on Monday and finishes on Sunday.
func setTimesLastWeek() {
	weekday := time.Now().Weekday()
	lastWeek := time.Now().AddDate(0, 0, int(-(7 + weekday)))
	lastWeekStop := lastWeek.AddDate(0, 0, 6)
	setTimes(lastWeek, lastWeekStop)
}

// Set the time range variables so that we look at data from this week
// where the week starts on Monday and finishes on Sunday.
func setTimesThisWeek() {
	weekday := time.Now().Weekday()
	lastWeek := time.Now().AddDate(0, 0, int(-(weekday)))
	lastWeekStop := lastWeek.AddDate(0, 0, 6)
	setTimes(lastWeek, lastWeekStop)
}

// Set the time range variables so that we look at data from this month only.
func setTimesThisMonth() {
	monthStart := time.Now().AddDate(0, 0, -(time.Now().Local().Day() - 1))
	year := time.Now().Year()
	month := time.Now().Local().Month() // Jan = 1
	days := 0
	if year%4 == 0 {
		days = daysInLeapMonth[month-1]
	} else {
		days = daysInMonth[month-1]
	}
	monthEnd := monthStart.AddDate(0, 0, days-1)
	setTimes(monthStart, monthEnd)
}

// Set the time range variables so that we look at data from the month
// prior to the current month.
func setTimesLastMonth() {
	monthStart := time.Now().AddDate(0, -1, -(time.Now().Local().Day() - 1))
	year := time.Now().Year()
	month := time.Now().Local().Month() // Jan = 1
	days := 0
	if year%4 == 0 {
		days = daysInLeapMonth[month-1]
	} else {
		days = daysInMonth[month-1]
	}
	monthEnd := monthStart.AddDate(0, 0, days-1)
	setTimes(monthStart, monthEnd)
}

// Read in the data file
func readFile(filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error opening %s: %v", filename, err)
		os.Exit(1)
	}
	defer fp.Close()

	fileScanner := bufio.NewScanner(fp)
	lineNo := 0
	for fileScanner.Scan() {
		lineNo++
		text := fileScanner.Text()
		if isInTimeRange(text) {
			// get hour of day, day of week, and day of month
			addTime(text)
			total++
		}
	}
	return nil
}

// Add the time information so we can report time of day, day of month.
// Example access.log entry:
// 192.168.0.22 - - [23/Aug/2022:01:07:47 -0400] "GET / HTTP/1.1" 200 2439 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"
// We want to pull just the date in brackets.
// TODO: Allow support of UTC instead of assuming local time.
func addTime(text string) {
	s := strings.SplitN(text, "[", 2)
	if len(s) != 2 {
		return // Invalid
	}
	s = strings.SplitN(s[1], "]", 2)
	if len(s) != 2 {
		return // Invalid
	}
	// s should look like: 23/Aug/2022:01:07:47 -0400
	const dateFormat = "02/Jan/2006:15:04:05 -0700"
	datetime, err := time.Parse(dateFormat, s[0])
	if err != nil {
		fmt.Printf("Invalid time string '%s': %v\n", s, err)
		return
	}

	// first get the day of month
	day := datetime.Local().Day() - 1

	// now get the month (Jan = 0)
	month := datetime.Local().Month() - 1

	// now get the year
	//year := datetime.Local().Year()

	/* now get the hour of day */
	hour := datetime.Local().Hour()

	// determine day of week (Sunday = 0)
	weekday := datetime.Local().Weekday()

	dow[weekday]++
	tod[hour]++
	moy[month]++
	dom[day]++
}

// Generate the output image
func generateOutput(fp *os.File, title string) {
	maxval := 0
	lastDay := 0
	width := 0
	height := 0

	// Display day of month data
	if outputType == DayOfMonth {
		/* get max value and last day we have data for */
		for loop := 0; loop < 31; loop++ {
			maxval = max(maxval, dom[loop])
			if dom[loop] > 0 {
				lastDay = loop
			}
		}
		maxval = calcMax(maxval)
		// Create graph output
		width = LeftPad + ((lastDay + 1) * DataWidth) + RightPad
		width = max(width, 500)
		height = TopPad + DataHeight + BottomPad
	} else if outputType == TimeOfDay {
		// get max value
		for loop := 0; loop < 24; loop++ {
			maxval = max(maxval, tod[loop])
		}
		maxval = calcMax(maxval)
		// Create graph output
		width = LeftPad + (24 * DataWidth) + RightPad
		height = TopPad + DataHeight + BottomPad
	} else if outputType == DayOfWeek {
		// get max value
		for loop := 0; loop < 7; loop++ {
			maxval = max(maxval, dow[loop])
		}
		maxval = calcMax(maxval)
		// Create graph output
		width = LeftPad + (7 * DataWidth * 3) + RightPad
		height = TopPad + DataHeight + BottomPad
	} else if outputType == MonthOfYear {
		// get max value
		for loop := 0; loop < 12; loop++ {
			maxval = max(maxval, moy[loop])
		}
		maxval = calcMax(maxval)
		// Create graph output
		width = LeftPad + (24 * DataWidth) + RightPad
		width = max(width, 500)
		height = TopPad + DataHeight + BottomPad
	}

	// allocate image
	im_out := ilibgo.CreateImage(width, height)
	/* first color is background color */
	black, _ := ilibgo.AllocNamedColor("black")
	grey, _ := ilibgo.AllocNamedColor("grey")
	red, _ := ilibgo.AllocNamedColor("darkred")
	blue, _ := ilibgo.AllocNamedColor("navy")
	gc := ilibgo.CreateGraphicsContext()
	ilibgo.SetForeground(&gc, grey)
	ilibgo.FillRectangle(im_out, gc, 0, 0, width, height)

	// draw black border
	ilibgo.SetForeground(&gc, black)
	ilibgo.DrawRectangle(im_out, gc, 0, 0, width-1, height-1)

	// draw title
	smallFont, _ := ilibgo.LoadFontFromData("utrg10", font.Font_UTRG__10())
	largeFont, _ := ilibgo.LoadFontFromData("utrg18", font.Font_UTRG__18())

	// Draw title
	ilibgo.SetFont(&gc, largeFont)
	titleW, titleH, _ := ilibgo.TextDimensions(gc, largeFont, title)
	// TODO: make sure titleW is not larger than TopPad if we allow user to
	// specify font
	ilibgo.SetForeground(&gc, grey)
	ilibgo.DrawString(im_out, gc, (width-titleW)/2+1, titleH+1, title)
	ilibgo.SetForeground(&gc, black)
	ilibgo.DrawString(im_out, gc, (width-titleW)/2, titleH, title)

	// write "Retrievals" to the left of the y axis
	ilibgo.SetFont(&gc, smallFont)
	ilibgo.DrawString(im_out, gc, 5, TopPad+DataHeight/2-5,
		yAxisLabel)

	// label the y axis  with the top value
	temp := abbreviatedNumber(maxval)
	w, _, _ := ilibgo.TextDimensions(gc, smallFont, temp)
	ilibgo.DrawString(im_out, gc, LeftPad-3-w, TopPad+3, temp)
	ilibgo.SetForeground(&gc, blue)
	ilibgo.DrawLine(im_out, gc, LeftPad-3, TopPad, LeftPad, TopPad)

	// draw x and y axis in blue
	ilibgo.DrawLine(im_out, gc, LeftPad, TopPad, LeftPad, TopPad+DataHeight)
	if graphType == LineGraph {
		ilibgo.DrawLine(im_out, gc, LeftPad, TopPad+DataHeight, width-RightPad,
			TopPad+DataHeight)
	} else {
		ilibgo.DrawLine(im_out, gc, LeftPad, TopPad+DataHeight,
			width-RightPad+(DataWidth/2),
			TopPad+DataHeight)
	}

	/* draw dashed lines horizontally */
	ilibgo.SetLineStyle(&gc, ilibgo.LineOnOffDash)
	for loop := 0; loop <= 4; loop++ {
		y := TopPad + ((DataHeight / 5) * loop)
		if graphType == LineGraph {
			ilibgo.DrawLine(im_out, gc, LeftPad, y, width-RightPad, y)
		} else {
			ilibgo.DrawLine(im_out, gc, LeftPad, y,
				width-RightPad+(DataWidth/2), y)
		}
	}
	ilibgo.SetLineStyle(&gc, ilibgo.LineSolid)

	if outputType == DayOfMonth {
		// label this as day of month on x axis
		ilibgo.SetForeground(&gc, black)
		label := "Day of the Month"
		axisW, _, _ := ilibgo.TextDimensions(gc, smallFont, label)
		x := LeftPad + ((width - RightPad - LeftPad) / 2) - (axisW / 2)
		ilibgo.DrawString(im_out, gc, x, height-35, label)
		// now draw a line for each point
		lastx := LeftPad + DataWidth
		lasty := TopPad + int((float64(maxval-dom[0])/float64(maxval))*float64(DataHeight))
		if graphType == BarGraph {
			ilibgo.SetForeground(&gc, red)
			ilibgo.FillRectangle(im_out, gc, lastx-(DataWidth/2)+1, lasty,
				DataWidth-2, (TopPad+DataHeight)-lasty)
		}
		// draw a tic mark
		ilibgo.SetForeground(&gc, blue)
		ilibgo.DrawLine(im_out, gc, lastx, TopPad+DataHeight, lastx,
			TopPad+DataHeight+3)
		ilibgo.SetForeground(&gc, black)
		ilibgo.DrawString(im_out, gc, lastx-2, TopPad+DataHeight+11,
			"1")
		for loop := 0; loop < lastDay; loop++ {
			x := LeftPad + ((loop + 1) * DataWidth)
			y := TopPad + int(((float64(maxval-dom[loop]) / float64(maxval)) * float64(DataHeight)))
			// draw the line graph
			ilibgo.SetForeground(&gc, red)
			if graphType == LineGraph {
				ilibgo.SetLineWidth(&gc, lineWidth)
				ilibgo.DrawLine(im_out, gc, lastx, lasty, x, y)
				ilibgo.SetLineWidth(&gc, 1)
			} else {
				ilibgo.FillRectangle(im_out, gc, x-(DataWidth/2)+1, y,
					DataWidth-2, (TopPad+DataHeight)-y)
			}
			// draw a tic mark
			ilibgo.SetForeground(&gc, blue)
			ilibgo.DrawLine(im_out, gc, x, TopPad+DataHeight, x,
				TopPad+DataHeight+3)
			// label the x axis
			ilibgo.SetForeground(&gc, black)
			temp := fmt.Sprintf("%d", (loop + 1))
			if loop < 9 {
				ilibgo.DrawString(im_out, gc, x-2, TopPad+DataHeight+11,
					temp)
			} else {
				ilibgo.DrawString(im_out, gc, x-5, TopPad+DataHeight+11,
					temp)
			}
			// save x and y for the next point
			lastx = x
			lasty = y
		}
	} else if outputType == TimeOfDay {
		// label this as time of day on x axis
		ilibgo.SetForeground(&gc, black)
		ilibgo.DrawString(im_out, gc, LeftPad+80, height-35,
			"Hour of the Day")
		// now draw a line for each point
		lastx := LeftPad + DataWidth
		lasty := TopPad + int((float64(maxval-tod[0])/float64(maxval))*float64(DataHeight))
		if graphType == BarGraph {
			ilibgo.SetForeground(&gc, red)
			ilibgo.FillRectangle(im_out, gc, lastx-(DataWidth/2)+1, lasty,
				DataWidth-2, (TopPad+DataHeight)-lasty)
		}
		// draw a tic mark
		ilibgo.SetForeground(&gc, blue)
		ilibgo.DrawLine(im_out, gc, lastx, TopPad+DataHeight, lastx,
			TopPad+DataHeight+3)
		ilibgo.SetForeground(&gc, black)
		ilibgo.DrawString(im_out, gc, lastx-2, TopPad+DataHeight+11, "1")
		for loop := 1; loop < 24; loop++ {
			x := LeftPad + ((loop + 1) * DataWidth)
			y := TopPad + int(((float64(maxval-tod[loop]) / float64(maxval)) * float64(DataHeight)))
			// draw the line graph
			ilibgo.SetForeground(&gc, red)
			if graphType == LineGraph {
				ilibgo.SetLineWidth(&gc, lineWidth)
				ilibgo.DrawLine(im_out, gc, lastx, lasty, x, y)
				ilibgo.SetLineWidth(&gc, 1)
			} else {
				ilibgo.FillRectangle(im_out, gc, x-(DataWidth/2)+1, y,
					DataWidth-2, (TopPad+DataHeight)-y)
			}
			// draw a tic mark
			ilibgo.SetForeground(&gc, blue)
			ilibgo.DrawLine(im_out, gc, x, TopPad+DataHeight, x,
				TopPad+DataHeight+3)
			// label the x axis
			temp := fmt.Sprintf("%d", (loop + 1))
			ilibgo.SetForeground(&gc, black)
			if loop < 9 {
				ilibgo.DrawString(im_out, gc, x-2, TopPad+DataHeight+11,
					temp)
			} else {
				ilibgo.DrawString(im_out, gc, x-5, TopPad+DataHeight+11,
					temp)
			}
			/* save x and y for the next point */
			lastx = x
			lasty = y
		}
	} else if outputType == DayOfWeek {
		// label this as time of day on x axis
		ilibgo.SetForeground(&gc, black)
		ilibgo.DrawString(im_out, gc, LeftPad+80, height-35,
			"Day of Week")
		// now draw a line for each point
		lastx := LeftPad + DataWidth*3
		lasty := TopPad + int((float64(maxval-tod[0])/float64(maxval))*float64(DataHeight))
		if graphType == BarGraph {
			ilibgo.SetForeground(&gc, red)
			ilibgo.FillRectangle(im_out, gc, lastx-(DataWidth/2)+1, lasty,
				DataWidth-2, (TopPad+DataHeight)-lasty)
		}
		// draw a tic mark
		ilibgo.SetForeground(&gc, blue)
		ilibgo.DrawLine(im_out, gc, lastx, TopPad+DataHeight, lastx,
			TopPad+DataHeight+3)
		ilibgo.SetForeground(&gc, black)
		ilibgo.DrawString(im_out, gc, lastx-2, TopPad+DataHeight+11,
			"Sun")
		for loop := 1; loop < 7; loop++ {
			x := LeftPad + ((loop + 1) * DataWidth * 3)
			y := TopPad + int((float64(maxval-dow[loop])/float64(maxval))*float64(DataHeight))
			// draw the line graph
			ilibgo.SetForeground(&gc, red)
			if graphType == LineGraph {
				ilibgo.SetLineWidth(&gc, lineWidth)
				ilibgo.DrawLine(im_out, gc, lastx, lasty, x, y)
				ilibgo.SetLineWidth(&gc, 1)
			} else {
				ilibgo.FillRectangle(im_out, gc, x-(DataWidth/2)+1, y,
					DataWidth-2, (TopPad+DataHeight)-y)
			}
			// draw a tic mark
			ilibgo.SetForeground(&gc, blue)
			ilibgo.DrawLine(im_out, gc, x, TopPad+DataHeight, x,
				TopPad+DataHeight+3)
			// label the x axis
			temp := wdays[loop]
			ilibgo.SetForeground(&gc, black)
			if loop < 9 {
				ilibgo.DrawString(im_out, gc, x-9, TopPad+DataHeight+11,
					temp)
			} else {
				ilibgo.DrawString(im_out, gc, x-12, TopPad+DataHeight+11,
					temp)
			}
			// save x and y for the next point
			lastx = x
			lasty = y
		}
	} else if outputType == MonthOfYear {
		// label this as time of day on x axis
		ilibgo.SetForeground(&gc, black)
		ilibgo.DrawString(im_out, gc, LeftPad+80, height-35,
			"Month of Year")
		// now draw a line for each point
		lastx := LeftPad + DataWidth*2
		lasty := TopPad + int((float64(maxval-moy[0])/float64(maxval))*float64(DataHeight))
		if graphType == BarGraph {
			ilibgo.SetForeground(&gc, red)
			ilibgo.FillRectangle(im_out, gc, lastx-(DataWidth/2)+1, lasty,
				DataWidth-2, (TopPad+DataHeight)-lasty)
		}
		// draw a tic mark
		ilibgo.SetForeground(&gc, blue)
		ilibgo.DrawLine(im_out, gc, lastx, TopPad+DataHeight, lastx,
			TopPad+DataHeight+3)
		ilibgo.SetForeground(&gc, black)
		ilibgo.DrawString(im_out, gc, lastx-2, TopPad+DataHeight+11,
			months[0])
		for loop := 1; loop < 12; loop++ {
			x := LeftPad + ((loop + 1) * DataWidth * 2)
			y := TopPad + (int)((float64(maxval-moy[loop])/float64(maxval))*float64(DataHeight))
			// draw the line graph
			ilibgo.SetForeground(&gc, red)
			if graphType == LineGraph {
				ilibgo.SetLineWidth(&gc, lineWidth)
				ilibgo.DrawLine(im_out, gc, lastx, lasty, x, y)
				ilibgo.SetLineWidth(&gc, 1)
			} else {
				ilibgo.FillRectangle(im_out, gc, x-(DataWidth/2)+1, y,
					DataWidth-2, (TopPad+DataHeight)-y)
			}
			// draw a tic mark
			ilibgo.SetForeground(&gc, blue)
			ilibgo.DrawLine(im_out, gc, x, TopPad+DataHeight, x,
				TopPad+DataHeight+3)
			// label the x axis
			ilibgo.SetForeground(&gc, black)
			if loop < 9 {
				ilibgo.DrawString(im_out, gc, x-9, TopPad+DataHeight+11,
					months[loop])
			} else {
				ilibgo.DrawString(im_out, gc, x-12, TopPad+DataHeight+11,
					months[loop])
			}
			// save x and y for the next point
			lastx = x
			lasty = y
		}
	}

	if displayHeader {
		temp := fmt.Sprintf("Total views for interval: %d", total)
		ilibgo.SetForeground(&gc, black)
		ilibgo.DrawString(im_out, gc, 5, height-13, temp)

		if len(prettyStart) > 0 && len(prettyStop) == 0 {
			temp = fmt.Sprintf("Time range: after %s", prettyStart)
		} else if len(prettyStop) > 0 && len(prettyStart) == 0 {
			temp = fmt.Sprintf("Time range: prior to %s", prettyStop)
		} else if len(prettyStart) > 0 && len(prettyStop) > 0 {
			temp = fmt.Sprintf("Time range: %s through %s",
				prettyStart, prettyStop)
		} else {
			temp = fmt.Sprintf("Time range: all")
		}
		ilibgo.DrawString(im_out, gc, 5, height-23, temp)
	}

	// Write output file.
	ilibgo.WriteImageFile(fp, im_out, ilibgo.FormatPNG)
}

// Calculate a nice round number to use as the maximum for the y axis.
// For example, if our max is 178, then use 200, etc.
func calcMax(curMax int) int {
	if curMax == 0 {
		return 1
	}

	curMaxF := float64(curMax)

	tens := int(math.Log10(curMaxF))
	nextRound := int(math.Pow(10.0, float64(tens+1)))
	divVal := nextRound / curMax
	switch divVal {
	case 10, 9, 8, 7, 6, 5:
		nextRound = nextRound / 5
	case 4:
		nextRound = nextRound / 4
	case 3:
		nextRound = int(0.4 * float64(nextRound))
	case 2:
		nextRound = nextRound / 2
	}

	retVal := max(5, nextRound)

	return retVal
}

// Check to see if the entry is within our time window
func isInTimeRange(entry string) bool {
	s := strings.SplitN(entry, "[", 2)
	if len(s) != 2 {
		return false
	}
	s = strings.SplitN(s[1], "]", 2)
	if (len(s)) != 2 {
		return false
	}
	timeStr := s[0]
	const dateFormat = "02/Jan/2006:15:04:05 -0700"
	datetime, err := time.Parse(dateFormat, timeStr)
	if err != nil {
		return false
	}

	if datetime.After(startTime) && datetime.Before(stopTime) {
		return true
	} else {
		return false
	}
}

// Convert a large number such as "1000000" into an abbreviated form like "1M".
func abbreviatedNumber(num int) string {
	if num >= 10000000 {
		num = num / 1000000
		return fmt.Sprintf("%dM", num)
	} else if num >= 1000000 {
		numF := float64(num) / 1000000.0
		return fmt.Sprintf("%.1fM", numF)
	} else if num >= 10000 {
		num = num / 1000
		return fmt.Sprintf("%dk", num)
	} else if num >= 1000 {
		numF := float64(num) / 1000.0
		return fmt.Sprintf("%.1fk", numF)
	}
	return fmt.Sprintf("%d", num)
}
