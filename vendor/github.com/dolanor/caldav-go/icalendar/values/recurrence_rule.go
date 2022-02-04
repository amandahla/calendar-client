package values

import (
	"fmt"
	"github.com/dolanor/caldav-go/icalendar/properties"
	"github.com/dolanor/caldav-go/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// The recurrence rule, if specified, is used in computing the recurrence set. The recurrence set is the complete set
// of recurrence instances for a calendar component. The recurrence set is generated by considering the initial
// "DTSTART" property along with the "RRULE", "RDATE", "EXDATE" and "EXRULE" properties contained within the iCalendar
// object. The "DTSTART" property defines the first instance in the recurrence set. Multiple instances of the "RRULE"
// and "EXRULE" properties can also be specified to define more sophisticated recurrence sets. The final recurrence
// set is generated by gathering all of the start date/times generated by any of the specified "RRULE" and "RDATE"
// properties, and excluding any start date/times which fall within the union of start date/times generated by any
// specified "EXRULE" and "EXDATE" properties. This implies that start date/times within exclusion related properties
// (i.e., "EXDATE" and "EXRULE") take precedence over those specified by inclusion properties
// (i.e., "RDATE" and "RRULE"). Where duplicate instances are generated by the "RRULE" and "RDATE" properties, only
// one recurrence is considered. Duplicate instances are ignored.

// The "DTSTART" and "DTEND" property pair or "DTSTART" and "DURATION" property pair, specified within the iCalendar
// object defines the first instance of the recurrence. When used with a recurrence rule, the "DTSTART" and "DTEND"
// properties MUST be specified in local time and the appropriate set of "VTIMEZONE" calendar components MUST be
// included. For detail on the usage of the "VTIMEZONE" calendar component, see the "VTIMEZONE" calendar component
// definition.

// Any duration associated with the iCalendar object applies to all members of the generated recurrence set. Any
// modified duration for specific recurrences MUST be explicitly specified using the "RDATE" property.
type RecurrenceRule struct {
	Frequency     RecurrenceFrequency
	Until         *DateTime
	Count         int
	Interval      int
	BySecond      []int
	ByMinute      []int
	ByHour        []int
	ByDay         []RecurrenceWeekday
	ByMonthDay    []int
	ByYearDay     []int
	ByWeekNumber  []int
	ByMonth       []int
	BySetPosition []int
	WeekStart     RecurrenceWeekday
}

var _ = log.Print

// the frequency an event recurs
type RecurrenceFrequency string

const (
	SecondRecurrenceFrequency RecurrenceFrequency = "SECONDLY"
	MinuteRecurrenceFrequency                     = "MINUTELY"
	HourRecurrenceFrequency                       = "HOURLY"
	DayRecurrenceFrequency                        = "DAILY"
	WeekRecurrenceFrequency                       = "WEEKLY"
	MonthRecurrenceFrequency                      = "MONTHLY"
	YearRecurrenceFrequency                       = "YEARLY"
)

// the frequency an event recurs
type RecurrenceWeekday string

const (
	MondayRecurrenceWeekday    RecurrenceWeekday = "MO"
	TuesdayRecurrenceWeekday                     = "TU"
	WednesdayRecurrenceWeekday                   = "WE"
	ThursdayRecurrenceWeekday                    = "TH"
	FridayRecurrenceWeekday                      = "FR"
	SaturdayRecurrenceWeekday                    = "SA"
	SundayRecurrenceWeekday                      = "SU"
)

// creates a new recurrence rule object for iCalendar
func NewRecurrenceRule(frequency RecurrenceFrequency) *RecurrenceRule {
	return &RecurrenceRule{Frequency: frequency}
}

var weekdayRegExp = regexp.MustCompile("MO|TU|WE|TH|FR|SA|SU")

// returns true if weekday is a valid constant
func (r RecurrenceWeekday) IsValidWeekDay() bool {
	return weekdayRegExp.MatchString(strings.ToUpper(string(r)))
}

var frequencyRegExp = regexp.MustCompile("SECONDLY|MINUTELY|HOURLY|DAILY|WEEKLY|MONTHLY|YEARLY")

// returns true if weekday is a valid constant
func (r RecurrenceFrequency) IsValidFrequency() bool {
	return frequencyRegExp.MatchString(strings.ToUpper(string(r)))
}

// returns the recurrence rule name for the iCalendar specification
func (r *RecurrenceRule) EncodeICalName() (properties.PropertyName, error) {
	return properties.RecurrenceRulePropertyName, nil
}

// encodes the recurrence rule value for the iCalendar specification
func (r *RecurrenceRule) EncodeICalValue() (string, error) {

	out := []string{fmt.Sprintf("FREQ=%s", strings.ToUpper(string(r.Frequency)))}

	if r.Until != nil {
		if encoded, err := r.Until.EncodeICalValue(); err != nil {
			return "", utils.NewError(r.EncodeICalValue, "unable to encode until date", r, err)
		} else {
			out = append(out, fmt.Sprintf("UNTIL=%s", encoded))
		}
	}

	if r.Count > 0 {
		out = append(out, fmt.Sprintf("COUNT=%d", r.Count))
	}

	if r.Interval > 0 {
		out = append(out, fmt.Sprintf("INTERVAL=%d", r.Interval))
	}

	if len(r.BySecond) > 0 {
		if encoded, err := intsToCSV(r.BySecond); err != nil {
			return "", utils.NewError(r.EncodeICalValue, "unable to encode by second value", r, err)
		} else {
			out = append(out, fmt.Sprintf("BYSECOND=%s", encoded))
		}
	}

	if len(r.ByMinute) > 0 {
		if encoded, err := intsToCSV(r.ByMinute); err != nil {
			return "", utils.NewError(r.EncodeICalValue, "unable to encode by minute value", r, err)
		} else {
			out = append(out, fmt.Sprintf("BYMINUTE=%s", encoded))
		}
	}

	if len(r.ByHour) > 0 {
		if encoded, err := intsToCSV(r.ByHour); err != nil {
			return "", utils.NewError(r.EncodeICalValue, "unable to encode by hour value", r, err)
		} else {
			out = append(out, fmt.Sprintf("BYHOUR=%s", encoded))
		}
	}

	if len(r.ByDay) > 0 {
		if encoded, err := daysToCSV(r.ByDay); err != nil {
			return "", utils.NewError(r.EncodeICalValue, "unable to encode by day value", r, err)
		} else {
			out = append(out, fmt.Sprintf("BYDAY=%s", encoded))
		}
	}

	if len(r.ByMonthDay) > 0 {
		if encoded, err := intsToCSV(r.ByMonthDay); err != nil {
			return "", utils.NewError(r.EncodeICalValue, "unable to encode by month day value", r, err)
		} else {
			out = append(out, fmt.Sprintf("BYMONTHDAY=%s", encoded))
		}
	}

	if len(r.ByYearDay) > 0 {
		if encoded, err := intsToCSV(r.ByYearDay); err != nil {
			return "", utils.NewError(r.EncodeICalValue, "unable to encode by year day value", r, err)
		} else {
			out = append(out, fmt.Sprintf("BYYEARDAY=%s", encoded))
		}
	}

	if len(r.ByWeekNumber) > 0 {
		if encoded, err := intsToCSV(r.ByWeekNumber); err != nil {
			return "", utils.NewError(r.EncodeICalValue, "unable to encode by week number value", r, err)
		} else {
			out = append(out, fmt.Sprintf("BYWEEKNO=%s", encoded))
		}
	}

	if len(r.ByMonth) > 0 {
		if encoded, err := intsToCSV(r.ByMonth); err != nil {
			return "", utils.NewError(r.EncodeICalValue, "unable to encode by month value", r, err)
		} else {
			out = append(out, fmt.Sprintf("BYMONTH=%s", encoded))
		}
	}

	if len(r.BySetPosition) > 0 {
		if encoded, err := intsToCSV(r.BySetPosition); err != nil {
			return "", utils.NewError(r.EncodeICalValue, "unable to encode by set position value", r, err)
		} else {
			out = append(out, fmt.Sprintf("BYSETPOS=%s", encoded))
		}
	}

	if r.WeekStart != "" {
		out = append(out, fmt.Sprintf("WKST=%s", r.WeekStart))
	}

	return strings.Join(out, ";"), nil
}

var rruleParamRegExp = regexp.MustCompile("(\\w+)\\s*=\\s*([^;]+)")

// decodes the recurrence rule value from the iCalendar specification
func (r *RecurrenceRule) DecodeICalValue(value string) error {

	matches := rruleParamRegExp.FindAllStringSubmatch(value, -1)
	if len(matches) <= 0 {
		return utils.NewError(r.DecodeICalValue, "no recurrence rules found", r, nil)
	}

	for _, match := range matches {
		if err := r.decodeICalValue(match[1], match[2]); err != nil {
			msg := fmt.Sprintf("unable to decode %s value", match[1])
			return utils.NewError(r.DecodeICalValue, msg, r, err)
		}
	}

	return nil

}

func (r *RecurrenceRule) decodeICalValue(name string, value string) error {

	switch name {
	case "FREQ":
		r.Frequency = RecurrenceFrequency(value)
	case "UNTIL":
		until := new(DateTime)
		if err := until.DecodeICalValue(value); err != nil {
			return utils.NewError(r.decodeICalValue, "invalid until value "+value, r, err)
		} else {
			r.Until = until
		}
	case "COUNT":
		if count, err := strconv.ParseInt(value, 10, 64); err != nil {
			return utils.NewError(r.decodeICalValue, "invalid count value "+value, r, err)
		} else {
			r.Count = int(count)
		}
	case "INTERVAL":
		if interval, err := strconv.ParseInt(value, 10, 64); err != nil {
			return utils.NewError(r.decodeICalValue, "invalid interval value "+value, r, err)
		} else {
			r.Interval = int(interval)
		}
	case "BYSECOND":
		if ints, err := csvToInts(value); err != nil {
			return utils.NewError(r.decodeICalValue, "invalid by second value "+value, r, err)
		} else {
			r.BySecond = ints
		}
	case "BYMINUTE":
		if ints, err := csvToInts(value); err != nil {
			return utils.NewError(r.decodeICalValue, "invalid by minute value "+value, r, err)
		} else {
			r.ByMinute = ints
		}
	case "BYHOUR":
		if ints, err := csvToInts(value); err != nil {
			return utils.NewError(r.decodeICalValue, "invalid by hour value "+value, r, err)
		} else {
			r.ByHour = ints
		}
	case "BYDAY":
		if days, err := csvToDays(value); err != nil {
			return utils.NewError(r.decodeICalValue, "invalid by day value "+value, r, err)
		} else {
			r.ByDay = days
		}
	case "BYMONTHDAY":
		if ints, err := csvToInts(value); err != nil {
			return utils.NewError(r.decodeICalValue, "invalid by month day value "+value, r, err)
		} else {
			r.ByMonthDay = ints
		}
	case "BYYEARDAY":
		if ints, err := csvToInts(value); err != nil {
			return utils.NewError(r.decodeICalValue, "invalid by year day value "+value, r, err)
		} else {
			r.ByYearDay = ints
		}
	case "BYWEEKNO":
		if ints, err := csvToInts(value); err != nil {
			return utils.NewError(r.decodeICalValue, "invalid by week number value "+value, r, err)
		} else {
			r.ByWeekNumber = ints
		}
	case "BYMONTH":
		if ints, err := csvToInts(value); err != nil {
			return utils.NewError(r.decodeICalValue, "unable to encode by month value "+value, r, err)
		} else {
			r.ByMonth = ints
		}
	case "BYSETPOS":
		if ints, err := csvToInts(value); err != nil {
			return utils.NewError(r.decodeICalValue, "unable to encode by set position value "+value, r, err)
		} else {
			r.BySetPosition = ints
		}
	case "WKST":
		r.WeekStart = RecurrenceWeekday(value)
	}

	return nil

}

// validates the recurrence rule value against the iCalendar specification
func (r *RecurrenceRule) ValidateICalValue() error {
	if !r.Frequency.IsValidFrequency() {
		return utils.NewError(r.ValidateICalValue, "a frequency is required in all recurrence rules", r, nil)
	} else if r.Until != nil && r.Count > 0 {
		return utils.NewError(r.ValidateICalValue, "until and count values are mutually exclusive", r, nil)
	} else if found, fine := intsInRange(r.BySecond, 59); !fine {
		msg := fmt.Sprintf("by second value of %d is out of bounds", found)
		return utils.NewError(r.ValidateICalValue, msg, r, nil)
	} else if found, fine := intsInRange(r.ByMinute, 59); !fine {
		msg := fmt.Sprintf("by minute value of %d is out of bounds", found)
		return utils.NewError(r.ValidateICalValue, msg, r, nil)
	} else if found, fine := intsInRange(r.ByHour, 23); !fine {
		msg := fmt.Sprintf("by hour value of %d is out of bounds", found)
		return utils.NewError(r.ValidateICalValue, msg, r, nil)
	} else if err := daysInRange(r.ByDay); err != nil {
		return utils.NewError(r.ValidateICalValue, "by day value not in range", r, err)
	} else if found, fine := intsInRange(r.ByMonthDay, 31); !fine {
		msg := fmt.Sprintf("by month day value of %d is out of bounds", found)
		return utils.NewError(r.ValidateICalValue, msg, r, nil)
	} else if found, fine := intsInRange(r.ByYearDay, 366); !fine {
		msg := fmt.Sprintf("by year day value of %d is out of bounds", found)
		return utils.NewError(r.ValidateICalValue, msg, r, nil)
	} else if found, fine := intsInRange(r.ByMonth, 12); !fine {
		msg := fmt.Sprintf("by month value of %d is out of bounds", found)
		return utils.NewError(r.ValidateICalValue, msg, r, nil)
	} else if found, fine := intsInRange(r.BySetPosition, 366); !fine {
		msg := fmt.Sprintf("by month value of %d is out of bounds", found)
		return utils.NewError(r.ValidateICalValue, msg, r, nil)
	} else if err := dayInRange(r.WeekStart); r.WeekStart != "" && err != nil {
		return utils.NewError(r.ValidateICalValue, "week start value not in range", r, err)
	} else {
		return nil
	}
}

func intsToCSV(ints []int) (string, error) {
	csv := new(CSV)
	for _, i := range ints {
		*csv = append(*csv, fmt.Sprintf("%d", i))
	}
	return csv.EncodeICalValue()
}

func csvToInts(value string) (ints []int, err error) {
	csv := new(CSV)
	if ierr := csv.DecodeICalValue(value); err != nil {
		err = utils.NewError(csvToInts, "unable to decode CSV value", value, ierr)
		return
	}
	for _, v := range *csv {
		if i, ierr := strconv.ParseInt(v, 10, 64); err != nil {
			err = utils.NewError(csvToInts, "unable to parse int value "+v, value, ierr)
			return
		} else {
			ints = append(ints, int(i))
		}
	}
	return
}

func intsInRange(ints []int, max int) (int, bool) {
	for _, i := range ints {
		if i < -max || i > max {
			return i, false
		}
	}
	return 0, true
}

func daysInRange(days []RecurrenceWeekday) error {
	for _, day := range days {
		if err := dayInRange(day); err != nil {
			msg := fmt.Sprintf("day value %s is not in range", day)
			return utils.NewError(dayInRange, msg, days, err)
		}
	}
	return nil
}

var dayRegExp = regexp.MustCompile("(\\d{1,2})?(\\w{2})")

func dayInRange(day RecurrenceWeekday) error {
	var ordinal, weekday string
	if matches := dayRegExp.FindAllStringSubmatch(string(day), -1); len(matches) <= 0 {
		msg := fmt.Sprintf("weekday value %s is not in valid format", day)
		return utils.NewError(dayInRange, msg, day, nil)
	} else if len(matches[0]) > 2 {
		ordinal = matches[0][1]
		weekday = matches[0][2]
	} else {
		weekday = matches[0][1]
	}
	if !RecurrenceWeekday(weekday).IsValidWeekDay() {
		msg := fmt.Sprintf("weekday value %s is not valid", weekday)
		return utils.NewError(dayInRange, msg, day, nil)
	} else if i, err := strconv.ParseInt(ordinal, 10, 64); ordinal != "" && err != nil {
		msg := fmt.Sprintf("weekday ordinal value %d is not valid", i)
		return utils.NewError(dayInRange, msg, day, err)
	} else if i < -53 || i > 53 {
		msg := fmt.Sprintf("weekday ordinal value %d is not in range", i)
		return utils.NewError(dayInRange, msg, day, nil)
	} else {
		return nil
	}
}

func daysToCSV(days []RecurrenceWeekday) (string, error) {
	csv := new(CSV)
	for _, day := range days {
		*csv = append(*csv, strings.ToUpper(string(day)))
	}
	return csv.EncodeICalValue()
}

func csvToDays(value string) (days []RecurrenceWeekday, err error) {
	csv := new(CSV)
	if ierr := csv.DecodeICalValue(value); err != nil {
		err = utils.NewError(csvToInts, "unable to decode CSV value", value, ierr)
		return
	}
	for _, v := range *csv {
		days = append(days, RecurrenceWeekday(v))
	}
	return
}
