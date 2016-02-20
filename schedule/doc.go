// Package schedule implements all the scheduling of jobs. The approach used to
// implement is the same as the http golang library style.
// We have schedulers and schedulerfuncs that will be executed.
// Schedulers are in a chain, this chain will be run and will execute the first
// scheduler, then this will execute the next scheduler and so on.
// The unit of execution is the schedulerfunc that receives the result and the job
package schedule
