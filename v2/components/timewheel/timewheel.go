package timewheel

import (
	"sync"
	"time"
)

var (
	tw *TimeWheel
)

func Init() {
	tw = New(10000, 10000)
}
func Instance() *TimeWheel {
	return tw
}

type Job struct {
	Key  string
	Id   string
	Time int64
	Data interface{}
}

type request struct {
	action int
	job    *Job
}

func New(reqLen, dispatchLen int) *TimeWheel {
	me := &TimeWheel{
		jobsByTime:   map[int64]map[string]*Job{},
		jobsById:     map[string]*Job{},
		reqs:         make(chan *request, reqLen),
		dispatchJobs: make(chan *Job, dispatchLen),
		closed:       make(chan struct{}),
		callbacks:    map[string]Callback{},
	}
	go me.mainloop()
	go me.dispatch()
	return me
}

type Callback func(*Job)

type TimeWheel struct {
	jobsByTime   map[int64]map[string]*Job
	jobsById     map[string]*Job
	reqs         chan *request
	dispatchJobs chan *Job
	closed       chan struct{}
	mutex        sync.RWMutex
	callbacks    map[string]Callback
}

func (me *TimeWheel) Close() {
	close(me.closed)
}

func (me *TimeWheel) Sub(key string, f Callback) {
	me.mutex.Lock()
	defer me.mutex.Unlock()
	me.callbacks[key] = f
}

func (me *TimeWheel) Unsub(key string) {
	me.mutex.Lock()
	defer me.mutex.Unlock()
	delete(me.callbacks, key)
}

func (me *TimeWheel) dispatch() {
	for {
		select {
		case job := <-me.dispatchJobs:
			me.mutex.RLock()
			f, ok := me.callbacks[job.Key]
			me.mutex.RUnlock()
			if ok {
				go f(job)
			}
		case <-me.closed:
			return
		}
	}
}

func (me *TimeWheel) Add(job *Job) {
	select {
	case <-me.closed:
		return
	case me.reqs <- &request{
		job: job,
	}:
	}
}

func (me *TimeWheel) Delete(key, id string) {
	select {
	case <-me.closed:
		return
	case me.reqs <- &request{
		action: 1,
		job: &Job{
			Id:  id,
			Key: key,
		},
	}:
	}
}

func (me *TimeWheel) mainloop() {
	tk := time.NewTicker(time.Millisecond * 600)
	lastTime := time.Now().Unix()
	expiredJobTimes := map[int64]struct{}{}
LOOP:
	for {
		select {
		case now := <-tk.C:
			for jobTime := range expiredJobTimes {
				for _, job := range me.jobsByTime[jobTime] {
					delete(me.jobsById, job.Id)
					me.dispatchJobs <- job
				}
				delete(me.jobsByTime, jobTime)
			}
			expiredJobTimes = map[int64]struct{}{}
			for jobTime := lastTime + 1; jobTime <= now.Unix(); jobTime++ {
				for _, job := range me.jobsByTime[jobTime] {
					delete(me.jobsById, job.Id)
					me.dispatchJobs <- job
				}
				delete(me.jobsByTime, jobTime)
			}
			if now.Unix() > lastTime {
				lastTime = now.Unix()
			}
		case req := <-me.reqs:
			mapKey := req.job.Key + ":" + req.job.Id
			if req.action == 0 {
				if job, ok := me.jobsById[mapKey]; ok {
					delete(me.jobsByTime[job.Time], mapKey)
				}
				me.jobsById[mapKey] = req.job
				if _, ok := me.jobsByTime[req.job.Time]; !ok {
					me.jobsByTime[req.job.Time] = map[string]*Job{}
				}
				me.jobsByTime[req.job.Time][mapKey] = req.job
				if req.job.Time <= lastTime {
					expiredJobTimes[req.job.Time] = struct{}{}
				}
			} else {
				if job, ok := me.jobsById[mapKey]; ok {
					delete(me.jobsById, mapKey)
					delete(me.jobsByTime[job.Time], mapKey)
				}
			}
		case <-me.closed:
			break LOOP

		}
	}
	tk.Stop()
}
