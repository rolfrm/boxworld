package main
import "fmt"
import "rand"
//import "math"
type BSPData interface{
	GetPosition() Vec3
	GetSize() Vec3
}

type SPData interface {
	GetPosition() Vec3
	GetSize() Vec3
}

type SPNode interface {
	//GetChildren() []SPNode
	//GetData() []SPData
	//GetPosition() Vec3
	//GetSize() Vec3
	Insert(SPData)
}

type QTNode struct {
	Data []SPData
	Split [8]*QTNode
	Position Vec3
	Size Vec3
	dimOfInterest int
	Parent *QTNode
}

func (self *QTNode) IsInside2(obj SPData)(output int){
	objPos := obj.GetPosition()
	objSize := obj.GetSize()
	objmin := objPos.Sub(objSize)
	objmax := objPos.Add(objSize)
	boxMin := self.Position
	boxMax := self.Position.Add(self.Size)
	lowMin := boxMin.Sub(objmin)
	lowMax := boxMin.Sub(objmax)
	highMin := boxMax.Sub(objmin)
	highMax := boxMax.Sub(objmax)
	fmt.Println(lowMin,lowMax,highMin,highMax)
	if (lowMax.Max() > 0) || ( highMin.Min() < 0) {
		fmt.Println("All outside!")
		return 0
	}
	if lowMin.Min() < 0 && highMax.Max() > 0 {
		fmt.Println("All inside")
		return 1
	}
	fmt.Println("both")
	return 2
}

func (self *QTNode) IsInside(obj SPData)(output int){
	var objPos Vec3 = obj.GetPosition()
	var objSize Vec3 = obj.GetSize()
	var boxPos Vec3 = self.Position.Add(self.Size.Scale(0.5))
	var boxSize Vec3= self.Size.Scale(0.5)
	var diff Vec3= objPos.Sub(boxPos).Abs()
	a := diff.Add(objSize).Sub(boxSize)
	if  a.Max() < 0 {
		return 1
	}
	a = diff.Sub(objSize).Sub(boxSize)
	if a.Max() < 0 {
		return 2
	}
	
	return 0
	
}


func (self *QTNode) SplitAt(mid Vec3)(out [8]*QTNode){
	for i:= 0; i < 8;i++ {
		out[i] = new(QTNode)
	}
	vMin := self.Position
	out[0].Position = vMin
	out[1].Position = Vec3{mid.X,vMin.Y,vMin.Z}
	out[2].Position = Vec3{vMin.X,mid.Y,vMin.Z}
	out[3].Position = Vec3{mid.X,mid.Y,vMin.Z}
	out[4].Position = Vec3{vMin.X,vMin.Y,mid.Z}
	out[5].Position = Vec3{mid.X,vMin.Y,mid.Z}
	out[6].Position = Vec3{vMin.X,mid.Y,mid.Z}
	out[7].Position = Vec3{mid.X,mid.Y,mid.Z}
	ls := mid.Sub(vMin)
	ms := self.Size.Sub(mid.Sub(vMin))
	out[0].Size = ls
	out[1].Size = Vec3{ms.X,ls.Y,ls.Z}
	out[2].Size = Vec3{ls.X,ms.Y,ls.Z}
	out[3].Size = Vec3{ms.X,ms.Y,ls.Z}
	out[4].Size = Vec3{ls.X,ls.Y,ms.Z}
	out[5].Size = Vec3{ms.X,ls.Y,ms.Z}
	out[6].Size = Vec3{ls.X,ms.Y,ms.Z}
	out[7].Size = Vec3{ms.X,ms.Y,ms.Z}
	//for i:= 0; i < 8; i++ {
		//fmt.Println("box:" , out[i].Position, " " , out[i].Size)
	//}
	return out
}



func (self *QTNode) Len() int{
	return len(self.Data)
}

func (self *QTNode) Less(i, j int)(bool){
	var diff Vec3 = self.Data[i].GetPosition().Sub(self.Data[j].GetPosition())
	return diff.GetComponent(self.dimOfInterest) < 0
}

func (self *QTNode) Swap(i,j int) {
	var  buf SPData = self.Data[i]
	self.Data[i] = self.Data[j]
	self.Data[j] = buf
}

func (self *QTNode) GetMidPoint()(out Vec3){
	out = Vec3{0,0,0}
	var n float32 = 0
	for i := 0; i < len(self.Data);i++ {
		out = out.Add(self.Data[i].GetPosition())
		n += 1
	}
	out = out.Scale(1/float32(n))
	return out
}

func octTreeCellTester(mid Vec3, obj SPData)(func(int)(bool)){
	var cellTester Vec3 = obj.GetPosition().Sub(mid).ElemDiv(obj.GetSize())
	//fmt.Println("cellTester", cellTester)
	return func(cell int)(bool) {
		switch cell{
		case 0: if cellTester.X > -1 && cellTester.Y > -1 && cellTester.Z > -1 { return true }
		case 1: if cellTester.X < 1 && cellTester.Y > -1 && cellTester.Z > -1 {return true }
		
		case 2: if cellTester.X > -1 && cellTester.Y < 1 && cellTester.Z > -1 {return true }
		case 3: if cellTester.X < 1 && cellTester.Y < 1 && cellTester.Z > -1 {return true }
		
		case 4: if cellTester.X > -1 && cellTester.Y > -1 && cellTester.Z < 1 {return true }
		case 5: if cellTester.X < 1 && cellTester.Y > -1 && cellTester.Z < 1 {return true }
		
		case 6: if cellTester.X > -1 && cellTester.Y < 1 && cellTester.Z < 1 {return true }
		case 7: if cellTester.X < 1 && cellTester.Y < 1 && cellTester.Z < 1 {return true }
		}
		return false
	}
}

func (self *QTNode) Divide() {
	//Lets assume that the midpoint is the optimal split point
	if len(self.Data) < 4 {
		//fmt.Println("own length:", len(self.Data))
		return
	}
	var mid Vec3 = self.GetMidPoint()
	var bins [8]*QTNode = self.SplitAt(mid)
	var newData []SPData
	var testBins [8] int = [8]int{0,0,0,0,0,0,0,0}
	for i := 0; i < len(self.Data); i++ {
		var inGen bool = false
		for j:= 0; j < 8; j++ {
			isInside := bins[j].IsInside(self.Data[i])
			if  isInside ==1{
				testBins[j] +=1
				bins[j].Insert(self.Data[i])	
			}else if isInside == 2 && inGen == false{
				inGen = true
				newData = append(newData,self.Data[i])
				break
			}


		}
	}
	//fmt.Println(testBins)
	self.Data = newData
	self.Split = bins
	//self.Position = mid
	for i:= 0; i < 8;i++ {
		//fmt.Println(i)
		self.Split[i].Divide()
		self.Split[i].Parent = self
	}
}

func (self *QTNode) Insert(ins SPData){
	self.Data = append(self.Data,ins)
}

func (self *QTNode) Traverse(depth int){
	
	//fmt.Println("lv: ", depth, " Imidiate children: ", len(self.Data))
	if self.Split[0] == nil {
		return
	}
	for i:= 0; i < len(self.Split);i++ {
		self.Split[i].Traverse(depth+1)
	}
}

func (self *QTNode) CountObjsInLeafNodes()(acc int){
	if self.Split[0] != nil {
		for i := 0 ; i < len(self.Split);i++ {
			acc += self.Split[i].CountObjsInLeafNodes()
		}
		return acc
	}
	return len(self.Data)*len(self.Data)
	
}


type BSPNode struct{
	Position Vec3
	Normal Vec3

	Split [2]*BSPNode
	Data []*BSPData
}

func (self *BSPNode) SplitInsert(obj *BSPData){

}

func (self *BSPNode) Insert(obj *BSPData){
	if self.Split[0] != nil {
		self.Data = append(self.Data,obj)
	}else{
		self.SplitInsert(obj)
	}
}

type testBox struct{
	Pos Vec3
	Size Vec3
}

func (self *testBox) GetPosition()Vec3{
	return self.Pos
}

func (self *testBox) GetSize() Vec3{
	return self.Size
}

func (self *QTNode) Find(obj SPData) (index int, found bool){
	for i:= 0; i < len(self.Data);i++ {
		if self.Data[i] == obj {
			return i,true
		}
	}
	return 0,false
}

func (self *QTNode) RemoveObj(obj SPData){
	index, found := self.Find(obj)
	if found {
		self.Data = append(self.Data[:index],self.Data[index+1:]...)
		
	}
}

func (self *QTNode) Update(){
	
	for i:= 0; i < len(self.Data);i++ {
		objState := self.IsInside(self.Data[i])
		if objState != 1 {
			fmt.Println("Object wrongly placed" , self.Data[i])
			if self.Parent != nil {
				self.Parent.EvalBox(self.Data[i])
			}
			self.RemoveObj(self.Data[i])
			i -= 1
		}
	}
	if self.Split[0] == nil {
		return
	}
	for i:=0; i < 8; i++ {
		self.Split[i].Update()
	}

}

func (self *QTNode) EvalBox(box SPData){
	if self.IsInside(box) == 1 {
		fmt.Println(box)
		if self.Split[0] != nil {
			for i:=0; i < 8; i++ {
				if self.Split[i].IsInside(box) == 1 {
					self.Split[i].EvalBox(box)
					
				}
			}
		}
		self.Insert(box)
	}else if self.Parent != nil{
		fmt.Println("This happened")
		self.Parent.EvalBox(box)
	}else{
		fmt.Println("This shouldent happen...", self.IsInside(box),self)
	}
}
var iter int = 0
func doCollisionTest(obj1, obj2 SPData){
	//fmt.Println(iter)
	iter += 1
}

func (self *QTNode) TestCollisionsWith(obj SPData){
	for i:= 0; i < len(self.Data); i++ {
		doCollisionTest(obj,self.Data[i])
	}

	if self.Split[0] != nil {
		for i:= 0; i < 8; i++ {
			if self.Split[i].IsInside(obj) > 0 {
				self.Split[i].TestCollisionsWith(obj)
			}
		}
	}
	
}


func (self *QTNode) TestCollisions(){
	for i := 0; i < len(self.Data); i++ {
		self.TestCollisionsWith(self.Data[i])
		for j := i+1 ; j < len(self.Data); j++ {
			doCollisionTest(self.Data[i],self.Data[j])
		}
	}
	if self.Split[0] == nil {
		return
	}
	for i:= 0; i < 8 ; i++ {
		self.Split[i].TestCollisions()
	}
}


func BSPTest(){
	var qtn *QTNode = new(QTNode)
	qtn.Position = Vec3{0,0,0}
	qtn.Size = Vec3{5,5,5}
	fmt.Println(qtn.IsInside2(&testBox{Vec3{7,7,7},Vec3{1,1,1}}))
	qtn.Position = Vec3{-10000,-10000,-10000}
	qtn.Size = Vec3{20000,20000,20000}
	
	//return

	var nValues int = 1000
	var tdat []testBox = make([]testBox, nValues)
	for i := 0; i < nValues; i++ {
		//var n float32 = float32(i)
		var x float32 = 1000*rand.Float32()*0
		var y float32 = 1000*rand.Float32()
		var z float32 = 1000*rand.Float32()*0
		tdat[i] = testBox{Vec3{x,y,z},Vec3{1,1,1}}
		
		qtn.Insert(&tdat[i])
	}
	qtn.Divide()
	
	
	qtn.Update()
	fmt.Println("top level data", len(qtn.Data))
	tdat[1].Pos = qtn.Split[7].Position
	tdat[1].Size = Vec3{1,1,1}
	qtn.Update()
	fmt.Println("top level data" , len(qtn.Data))
	//qtn.TestCollisionsWith(&tdat[1])
	qtn.TestCollisions()
	fmt.Println("Counting...");
	fmt.Println("Collision detection: ", iter, " VS ", nValues*nValues)
}

