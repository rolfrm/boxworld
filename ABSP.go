package main
import "fmt"
func SPMap(f func(SPData), set []SPData){
	for i:=0; i < len(set);i++ {
		f(set[i])
	}
}

type ABSPNode struct{
	Data []SPData
	Split [2]*ABSPNode
	splitDim int
	splitPos float32
	IsSplit bool
	Root *ABSPNode
}

func (self *ABSPNode) Insert(obj SPData ){
	self.Data = append(self.Data,obj)
}


func (self *ABSPNode) GetMean()(Vec3){
	var µ Vec3 = Vec3{0,0,0}
	var n float32 = 0
	SPMap(func(obj SPData){
		µ = µ.Add(obj.GetPosition())
		n += 1
		
	},self.Data)
	return µ.Scale(1/n)
}

func (self *ABSPNode) GetVariance()(Vec3){
	var out Vec3 = Vec3{0,0,0}
	var mean Vec3 = self.GetMean()
	SPMap(func(obj SPData){
		out = out.Add(obj.GetPosition().Sub(mean).ElmPow(2))
	},self.Data)
	return out
}

func (self *ABSPNode) Divide(){
	
	if len(self.Data) < 1 {
		return
	}
	var mean Vec3 = self.GetMean()
	var variance Vec3 = self.GetVariance()
	var comp int = variance.BiggestComponent()
	var test *ABSPNode = new(ABSPNode)
	test.splitDim = comp
	test.splitPos = mean.GetComponent(comp)
	var caseCounter [3]int;
	SPMap(func(obj SPData){
		caseCounter[test.Cell(obj)] +=1
	},self.Data)
	//fmt.Println(caseCounter)
	if caseCounter[0] + caseCounter[1] > caseCounter[2] {
		self.IsSplit = true
		self.splitDim = test.splitDim
		self.splitPos = test.splitPos
		
		self.Split[0] = new(ABSPNode)
		self.Split[1] = new(ABSPNode)
		self.Split[0].Root = self.Root
		self.Split[1].Root = self.Root
		
		var newData []SPData
		SPMap(func( obj SPData){
			cell := self.Cell(obj)
			if cell != 2 {
				self.Split[cell].Insert(obj)
			}else{
				newData = append(newData, obj)
			}
		},self.Data)
		self.Data = newData
		self.Split[0].Divide()
		self.Split[1].Divide()
	}
}


func Cell(pos Vec3, size Vec3, splitPos float32, splitDim int) int{
	var oSize float32
	var diff float32

	if splitDim == 0 {
		oSize = size.X
		diff = pos.X
	}else if splitDim == 1 {
		oSize = size.Y
		diff = pos.Y
	}else{
		oSize = size.Z
		diff = pos.Z
	}
	diff -= splitPos

	
	if diff > oSize {
		return 1
	}

	if diff < -oSize {
		return 0
	}
	return 2
} 

func (self *ABSPNode) Cell(obj SPData) int{
	var size Vec3 = obj.GetSize()
	var pos Vec3 = obj.GetPosition()
	var oSize float32 = 0
	var diff float32 = 0

	if self.splitDim == 0 {
		oSize = size.X
		diff = pos.X
	}else if self.splitDim == 1 {
		oSize = size.Y
		diff = pos.Y
	}else{
		oSize = size.Z
		diff = pos.Z
	}
	diff -= self.splitPos

	
	if diff > oSize {
		return 1
	}

	if diff < -oSize {
		return 0
	}
	return 2
}

func (self *ABSPNode) findContainingChild(obj SPData)(*ABSPNode){
	if self.IsSplit == false {
		return self
	}

	cell := self.Cell(obj)
	if cell == 2 {
		return self
	}
	return self.Split[cell].findContainingChild(obj)

}

func (self *ABSPNode) Find(obj SPData)(index int){
	for i:= 0; i < len(self.Data); i++ {
		if self.Data[i] == obj {
			return i
		}
	}
	return -1
}

func (self *ABSPNode) RemoveObj(obj SPData){
	objIndex := self.Find(obj)
	if objIndex != -1 {
		self.Data = append(self.Data[:objIndex],self.Data[objIndex+1:]...)
	}
}

func (self *ABSPNode) Update(){
	for i:= 0; i < len(self.Data) ; i++ {
		var containingNode *ABSPNode  = self.Root.findContainingChild(self.Data[i])
		
		if containingNode == self {
			continue
		}
		containingNode.Insert(self.Data[i])
		self.Data = append(self.Data[:i], self.Data[i+1:]...)

		i -= 1
	}
	if self.IsSplit {
		self.Split[0].Update()
		self.Split[1].Update()
	}
}

func (self *ABSPNode) Traverse(i int){
	fmt.Println("recursion level: " , i, len(self.Data))
	if self.IsSplit {
		self.Split[0].Traverse(i+1)
		self.Split[1].Traverse(i+1)
	}
}

func (self *ABSPNode) cd(fcd func(o1,o2 SPData)){
	for i:= 0; i < len(self.Data); i++ {
		self.cdo(fcd,self.Data[i])
	}
	if self.IsSplit {
		self.Split[0].cd(fcd)
		self.Split[1].cd(fcd)
	}
}

func (self *ABSPNode) cdo(fcd func(o1,o2 SPData),obj SPData){
	data := self.Data
	for i:= 0; i < len(data);i++ {
		fcd(obj, data[i])
	}
	
	if self.IsSplit {
		cell := Cell(obj.GetPosition(),obj.GetSize(),self.splitPos,self.splitDim)
		if cell < 2 {
			
			self.Split[cell].cdo(fcd,obj)
		}else{
			self.Split[1].cdo(fcd,obj)
			self.Split[0].cdo(fcd,obj)
		}
	}
}


func (self *ABSPNode) CountObjects() int {
	if self.IsSplit {
		return len(self.Data) + self.Split[0].CountObjects() + self.Split[1].CountObjects()
	}
	return len(self.Data)
}



func ABSPTest(){

	absp := new(ABSPNode)
	absp.Root = absp
	var datalist []*testBox
	for i:= float32(0); i < 10000; i++{
		datalist = append(datalist,&testBox{Vec3{i/5,i/7,i/4},Vec3{1,1,1}})
	}
	for i:= 0; i < len(datalist);i++ {
		absp.Insert(datalist[i])
	}

	fmt.Println(absp.GetMean())
	fmt.Println(absp.GetVariance())
	absp.Divide()
	datalist[10].Pos =Vec3{0,0,0}
	datalist[15].Pos = Vec3{100,10,0}
	absp.Update()
	absp.Update()
	i  := 0
	/*absp.RunCollisionFunction(func(obj1, obj2 SPData){
		i += 1
		
	},nil)*/
	absp.cd(func(obj1,boj2 SPData){
		i +=  1
	})
	absp.Traverse(0)
	fmt.Println(i)
	fmt.Println(absp.CountObjects())
}