// Package schedule implements all the scheduling of jobs. The approach used to
// implement is the same as the http golang library style.
// We have Schedulers that will be executed (Run method).
// We have SchedulerFunc that is a wrap of a function to convert it to an schedule
// object and be able to run.
// like http middlewares, we create schedulers with functions that receive other schedulers as
// parameters and create a chain (schedule/chain)
package schedule
