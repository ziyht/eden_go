package etimer

import (
	"context"
	"fmt"
	"testing"
	"time"
)


func TestXxx(t *testing.T) {

	tm  := NewTimer()
	ctx := context.Background() 

	j := tm.AddInterval(ctx, time.Second * 2, func(j *Job)error{
		fmt.Printf("running start at %s\n", j.State())
		time.Sleep(time.Second * 3)
		fmt.Printf("running over at %s\n", time.Now())
		return nil
	})

	time.Sleep(time.Second * 300)
	fmt.Printf("%+v\n", j.State())



}