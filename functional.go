package main
import "container/list"
import "fmt"

func SliceMap(f func( interface{}),slc interface{}){
	nslc := slc.([] interface{})
	for i := 0;i < len(nslc); i++ {
		f(nslc[i])
	}
}

func FuncMap(f func( interface {} ) interface{}, lst *list.List) *list.List{
	output := list.New()
	for x := lst.Front(); x != nil; x= x.Next() {
		output.PushBack(f(x.Value.(int)))
	}
	return output;
}

func RecursiveFuncMap(f func( interface{} )(interface {}, *list.List), lst *list.List) *list.List{
	output := list.New()
	for x:=lst.Front(); x != nil; x = x.Next() {
		fout, rec := f(x.Value)
		output.PushBack(fout)
		if rec!= nil && rec.Len() > 0 {
			output.PushBack(RecursiveFuncMap(f,rec))
		}
	}
	return output
}

func ClosureTest(clo float32)(func (float32)(float32),func(float32)(float32)){
	return func(x float32)(float32){
		clo = clo*x 
		return clo
	}, func(x float32)(float32){
		clo = clo/x
		return clo
	}
	
}

func FunctionalTest(){
	data := list.New()
	
	data.PushBack(1)
	data.PushBack(2)
	data.PushBack(3)
	data.PushBack(4)
	nlist := FuncMap(func(x interface{}) interface{} {
	return x.(int) + 1
	},data)
	fmt.Println(nlist)
	
	nlist = RecursiveFuncMap(
		func(x interface{})(y interface{},outlist *list.List){
		outlist = list.New()
		for i:=0; i < x.(int);i++ {
			outlist.PushBack(i)
		}
		y = x.(int)
		return
	},data)

	nlist = RecursiveFuncMap(
		func(x interface{})(y interface{}, outlist * list.List){
		v,ok := x.(int)
		if ok {
			fmt.Println(v)
			return v,nil
		}
		v2,ok := x.(*list.List)
		if ok {
			fmt.Println("len: ", v2.Len())
			return 0,v2
		}

		return 0,nil

	},nlist)


}
