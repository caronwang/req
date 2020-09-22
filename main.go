package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Info struct {
	Idx int
	Pidx []int
}

func Print(a []*Info){
	for _,d := range a{
		fmt.Println(d.Idx," ",d.Pidx)
	}
}


var (
	Que_In = make(chan *Info)
	Que_Out = make(chan *Info)
)


func Check(a []*Info)  (r [][]*Info,err error) {

	res := make([]bool,len(a))
	b := a

	for len(b) != 0 {
		//fmt.Println("----")
		change := false
		tmp := []*Info{}
		doneArr := []*Info{}
		for _,d := range b{
			done := true
			for _,pid := range d.Pidx{
				if !res[pid] {
					tmp = append(tmp,d)
					//fmt.Println("failed!",d.Idx,d.Pidx)
					done = false
					break
				}
			}
			if !done{
				continue
			}
			//fmt.Println("ok --> ",d.Idx,d.Pidx)
			doneArr = append(doneArr,d)
			//res[d.Idx] = true
			change = true
		}

		r = append(r,doneArr)

		for _,d := range doneArr{
			//fmt.Println("set - ",d.Idx)
			res[d.Idx] = true
		}

		b = tmp
		if !change{
			Print(tmp)
			err = errors.New("error")
			return
		}
	}
	return
}


func Run(a [][]*Info){
	for i,d := range a {
		fmt.Println("start stage ",i+1)
		taskList := make(map[int]struct{})

		for _,item := range d {
			taskList[item.Idx] = struct{}{}
		}

		go func() {
			for _,item := range d {
				fmt.Println("\tsend --> ",item.Idx,item.Pidx)
				taskList[item.Idx] = struct{}{}
				Que_In <-item
			}
		}()

		if WaitForFinish(taskList) {
			fmt.Println("complete stage ",i+1)
		}else{
			fmt.Println("stage ",i+1,"failed")
			return
		}

	}

}


func WaitForFinish(list map[int]struct{}) bool{
	start := time.Now()
	t := time.NewTimer(time.Second*30)
	for {
		select {
		case data := <-Que_Out:
			if _,ok:= list[data.Idx];ok{
				delete(list,data.Idx )
			}
			if len(list) == 0{
				fmt.Println(" stage耗时",time.Since(start))
				return true
			}
		case <-t.C:
				fmt.Println("超时失败")
				return false
		}
	}
}



func main(){
	var a []*Info
	for i:=0;i<10;i++{
		a = append(a,&Info{Idx: i})
	}
	a[2].Pidx = append(a[2].Pidx,[]int{1,3}...)
	a[8].Pidx = append(a[8].Pidx,[]int{5,6}...)
	a[4].Pidx = append(a[4].Pidx,[]int{1,3}...)
	a[5].Pidx = append(a[5].Pidx,6)
	a[6].Pidx = append(a[6].Pidx,2)
	a[7].Pidx = append(a[7].Pidx,5)
	arr ,err  := Check(a)
	if err!=nil{
		log.Fatal(err)
	}

	go func() {
		for data:= range Que_In{
			//fmt.Println("start handle -->",data.Idx,data.Pidx)
			t:=time.Now()
			rand.Seed(time.Now().UnixNano())
			time.Sleep( time.Duration(rand.Intn(3))*time.Second)
			fmt.Println(" handle (",time.Since(t),") --> ",data.Idx,data.Pidx)
			Que_Out<-data
		}
	}()

	Run(arr)
}
