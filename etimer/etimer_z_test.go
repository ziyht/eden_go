package etimer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInterval(t *testing.T) {

	tm  := NewTimer(TimerOptions{Interval: time.Second/1000})
	ctx := context.Background() 
	cnt := 0
	j := tm.AddInterval(ctx, time.Second/100, func(j *Job)error{
		fmt.Printf("running at %s\n", j.State())
		cnt += 1
		return nil
	})

	time.Sleep(time.Second * 5)
	fmt.Printf("%d\n", cnt)
	fmt.Printf("%+v\n", j.State())
}

func TestLimit(t *testing.T) {

	tm := NewTimer(TimerOptions{Interval: time.Second/1000})

	j := NewJob(&JobOpts{
		Interval: time.Second/10,
		Times   : 30,
		CB      : func(j *Job)error{
			fmt.Printf("running at %s\n", j.State())
			return nil
		},
	})

	tm.AddJob(j)

	time.Sleep(time.Second * 5)
	fmt.Printf("%+v\n", j.State())
	assert.Equal(t, uint64(30), j.State().Runnings())
}

func TestCron(t *testing.T) {

	tm  := NewTimer()
	ctx := context.Background() 

	j, _ := tm.AddCron(ctx, "* * * * *", func(j *Job)error{
		fmt.Printf("running at %s\n", j.State())
		return nil
	})

	time.Sleep(time.Second * 100)
	fmt.Printf("%+v\n", j.State())
}

func TestCron2(t *testing.T) {

	tm := NewTimer(TimerOptions{Interval: time.Second/100})
	ctx := context.Background() 

	j1, _ := tm.AddCron(ctx, "@every 2s", func(j *Job)error{
		fmt.Printf("job1\n")
		return nil
	})

	j2, _ := tm.AddCron(ctx, "@every 1s", func(j *Job)error{
		fmt.Printf("job2\n")
		return nil
	})

	time.Sleep(time.Second * 10)
	fmt.Printf("%+v\n", j1.State())
	fmt.Printf("%+v\n", j2.State())
}