package NewTime

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"time"
)

type NewTime interface {
	getRegexpCompiled() (*regexp.Regexp, error)
	fillTimeFromString(string) error
}

type TimeISO8601 struct {
	Time time.Time
}

type TimeUnix struct {
	Time time.Time
}

func (t *TimeUnix) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Unix())
}

func (t *TimeISO8601) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format("2006-01-02 15:04:05.000"))
}

func (t *TimeUnix) UnmarshalJSON(b []byte) error {
	err := unmarshal(t, b)

	if err != nil {
		return err
	}

	return nil
}

func (t *TimeISO8601) UnmarshalJSON(b []byte) error {
	err := unmarshal(t, b)

	if err != nil {
		return err
	}

	return nil
}

func unmarshal(t NewTime, b []byte) error {
	err := validateJson(t, b)
	if err != nil {
		return err
	}

	str, err := getDateFromJson(t, b)

	if err != nil {
		return err
	}

	err = t.fillTimeFromString(str)

	if err != nil {
		return err
	}

	return nil
}

func validateJson(t NewTime, b []byte) error {
	r, err := t.getRegexpCompiled()
	if err != nil {

		return err
	}

	if !r.Match(b) {
		return errors.New("invalid time format")
	}

	return nil
}

func (t *TimeISO8601) getRegexpCompiled() (*regexp.Regexp, error) {
	return regexp.Compile(`\d{4}-(0[1-9]|1[0-2])-(0[1-9]|1[0-9]|2[0-9]|3[01]) (0[0-9]|1[0-9]|2[0-3]):([0-5]\d):([0-5]\d)\.\d{3}`)
}

func (t *TimeUnix) getRegexpCompiled() (*regexp.Regexp, error) {
	return regexp.Compile(`[\d]*`)
}

func getDateFromJson(t NewTime, b []byte) (string, error) {
	r, err := t.getRegexpCompiled()

	if err != nil {
		return "", err
	}

	return r.FindString(string(b)), nil
}

func (t *TimeISO8601) fillTimeFromString(str string) error {
	y, err := strconv.Atoi(str[:4])
	if err != nil {
		return err
	}

	m, err := strconv.Atoi(str[5:7])
	if err != nil {
		return err
	}

	d, err := strconv.Atoi(str[8:10])
	if err != nil {
		return err
	}

	hour, err := strconv.Atoi(str[11:13])
	if err != nil {
		return err
	}

	mins, err := strconv.Atoi(str[14:16])
	if err != nil {
		return err
	}

	sec, err := strconv.Atoi(str[17:19])
	if err != nil {
		return err
	}

	nsec, err := strconv.Atoi(str[20:])
	if err != nil {
		return err
	}

	t.Time = time.Date(y, time.Month(m), d, hour, mins, sec, nsec, time.UTC)

	return nil
}

func (t *TimeUnix) fillTimeFromString(str string) error {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(i, 0)

	return nil
}
