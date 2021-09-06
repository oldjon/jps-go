package main

import "fmt"

func main() {

	m := [][]byte{}
	m = append([][]byte{{0, 0, 0, 1, 0, 0, 0, 1, 0, 0}}, m...) //9
	m = append([][]byte{{0, 0, 0, 1, 0, 0, 0, 1, 0, 0}}, m...) //8
	m = append([][]byte{{0, 0, 0, 1, 0, 1, 0, 1, 0, 0}}, m...) //7
	m = append([][]byte{{0, 0, 0, 1, 0, 1, 0, 1, 0, 0}}, m...) //6
	m = append([][]byte{{0, 1, 0, 1, 0, 1, 0, 0, 0, 0}}, m...) //5
	m = append([][]byte{{0, 0, 0, 1, 0, 1, 1, 1, 1, 1}}, m...) //4
	m = append([][]byte{{0, 1, 0, 1, 0, 0, 0, 0, 0, 0}}, m...) //3
	m = append([][]byte{{0, 1, 0, 0, 1, 1, 1, 1, 0, 0}}, m...) //2
	m = append([][]byte{{0, 1, 0, 0, 0, 0, 0, 0, 0, 0}}, m...) //1
	m = append([][]byte{{0, 1, 0, 0, 1, 0, 0, 0, 0, 0}}, m...) //0
	//                   0  1  2  3  4  5  6  7  8  9

	jps := JPS{}
	jps.Init(m)
	err := jps.FindPath(&Pos{X: 0, Y: 0}, &Pos{X: 9, Y: 9})
	if err != nil {
		fmt.Println(err.Error())
	}
	jps.PrintPath()

	jps2 := JPS2{}
	jps2.Init(m)
	err = jps2.FindPath(&Pos{X: 0, Y: 0}, &Pos{X: 9, Y: 9})
	if err != nil {
		fmt.Println(err.Error())
	}
	jps2.PrintPath()
}
