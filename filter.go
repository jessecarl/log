package log

// Filter is a function that is used to manipulate the Data passed to a Logger.
// This could be adding fields, removing fields, or even setting the Data to
// nil. Filters need to be able to accept nil Data, as that is the most
// effective way of stopping Data from being logged.
type Filter func(lvl, threshold Level, data Data) Data
