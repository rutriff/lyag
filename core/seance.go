package core

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Seance struct {
	StaffId string
	Time    time.Time
	Service int
}

func (s Seance) String() string {
	return s.GetKey()
}

func (s *Seance) GetKey() string {
	return s.StaffId + "|" + s.Time.Format(time.RFC3339)
}

func ParseKey(key string) (int, time.Time, error) {
	parts := strings.Split(key, "|")

	if len(parts) != 2 {
		return 0, time.Time{}, errors.New("can't split key, expected two parts")
	}

	staffId, err := strconv.Atoi(parts[0])

	if err != nil {
		return 0, time.Time{}, err
	}

	t, err := time.Parse(time.RFC3339, parts[1])

	if err != nil {
		return staffId, time.Time{}, err
	}

	return staffId, t, nil
}
