package service

import "time"

func getLastDayOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}

func getLastDayOfWeek(date time.Time) time.Time {
	dayOrderNum := date.Weekday()
	// Если воскресенье, то уже конец недели, возвращаем воскресенье.
	if dayOrderNum == time.Sunday {
		return date
	}
	// Посчитать сколько осталось дней до воскресенья, т.к. воскресенье конец недели
	const numOfDaysInWeek = 7
	daysUntilSunday := numOfDaysInWeek - int(dayOrderNum)
	sunday := date.AddDate(0, 0, daysUntilSunday)
	return sunday
}
