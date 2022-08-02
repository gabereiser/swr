/*  Star Wars Role-Playing Mud
 *  Copyright (C) 2022 @{See Authors}
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */
package swr

import (
	"log"
	"sync"
	"time"
)

type ScheduledFunction struct {
	Repeat  bool
	Func    func()
	Seconds uint
	Current uint
}

func ScheduleFunc(fn func(), repeat bool, time uint) {
	sf := ScheduledFunction{
		Repeat:  repeat,
		Func:    fn,
		Seconds: time,
		Current: 0,
	}
	Scheduler().Schedule(&sf)
}

type SchedulerService struct {
	m     *sync.Mutex
	t     *time.Ticker
	funcs []*ScheduledFunction
	bt    time.Time
}

var _scheduler *SchedulerService

func Scheduler() *SchedulerService {
	if _scheduler == nil {
		log.Println("Starting Scheduler.")
		now := time.Now().UTC()
		delta := now.Sub(now.Truncate(time.Second))
		time.Sleep(time.Second - delta)
		_scheduler = &SchedulerService{
			t:     time.NewTicker(time.Duration(1) * time.Second),
			m:     &sync.Mutex{},
			funcs: []*ScheduledFunction{},
		}
		go func() {
			for {
				<-_scheduler.t.C
				now := time.Now().UTC()
				delta := now.Sub(now.Truncate(time.Second))
				time.Sleep(time.Second - delta)
				_scheduler.tick(time.Now().UTC())
			}
		}()
		_scheduler.bt = now
		log.Println("Scheduler Started.")
	}
	return _scheduler
}
func (s *SchedulerService) Lock() {
	s.m.Lock()
}

func (s *SchedulerService) Unlock() {
	s.m.Unlock()
}

func (s *SchedulerService) Schedule(function *ScheduledFunction) {
	s.Lock()
	defer s.Unlock()
	s.funcs = append(s.funcs, function)
}

func (s *SchedulerService) Remove(function *ScheduledFunction) {
	s.Lock()
	defer s.Unlock()
	ret := make([]*ScheduledFunction, 0)
	for _, f := range s.funcs {
		if f == function {
			continue
		}
		ret = append(ret, f)
	}
	s.funcs = ret
}
func (s *SchedulerService) tick(t time.Time) {

	removal := []*ScheduledFunction{}

	for _, fn := range s.funcs {
		fn.Current++
		if fn.Current == fn.Seconds {
			fn.Func()
			if !fn.Repeat {
				removal = append(removal, fn)
			}
			fn.Current = 0
		}
	}

	for _, fn := range removal {
		s.Remove(fn)
	}
}
