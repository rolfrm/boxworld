package main
import "glfw"
import "gl"
import "container/list"
import "fmt"
import "math"
import "time"

type GameObj struct {
	Pos Vec3
	Size Vec3
	Color Vec3
	Velocity Vec3
	Mass float32
	Friction float32
	Anim func(self *Box3D,t float32)
	Children *list.List
	Constraints []Constraint
	Rotation Vec3
	Body PhysicsBody
	IsGhost bool
}

func (self *GameObj) GetPosition()(Vec3){
	return self.Pos
}

func (self * GameObj) SetPosition(in Vec3){
	self.Pos = in
}

func (self * GameObj) GetFriction() float32 {
	return self.Friction
}

func (self *GameObj) GetVelocity()(Vec3){
	return self.Velocity
}

func (self *GameObj) SetVelocity(in Vec3)(){
	self.Velocity = in
}

func (self *GameObj) GetMass()(float32){
	return self.Mass
}

func (self * GameObj) SetMass(in float32){
	self.Mass = in
}


func (self *GameObj) GetBoundingBox()(bb Box){
	bb.Pos = self.Pos
	bb.Size = self.Size
	return
}

func (s *GameObj) GetBox3D()(output Box3D){
	output = Box3D{s.Pos,s.Size,s.Color,s.Anim}
	return
}

func (s *GameObj) GetConstraints()(out []Constraint){
	return s.Constraints
}

func (s *GameObj) GetChildren()(*list.Element){
	return s.Children.Front()
}

func (s *GameObj) AddChildren( childrenArray ... *GameObj){
	for it := 0; it < len(childrenArray);it++ {
		s.Children.PushBack(childrenArray[it])
	}
}

func (s *GameObj) GetSize()Vec3{
	return s.Size
}


func (self *GameObj) GetSubs() ([]PhysicsObject){
	
	output := []PhysicsObject{self}
	for i := self.Children.Front(); i != nil;i = i.Next() {
		child, ok := i.Value.(*GameObj)
		if ok {
			output = append(output,  child.GetSubs()...)
		}
	}
	return output
	
}

func (self *GameObj) GetBody()(*PhysicsBody){
	return &self.Body
}

func (self *GameObj) UpdatePhysics(){
	self.Pos = self.Body.Pos
	self.Velocity = self.Body.Velocity
}



type World struct{
	GameObjects *list.List
	GameObjectTree *ABSPNode
}

func (self *World) Init(){
	self.GameObjects = new(list.List)
}

func (self *World) Add(obj * GameObj){
	self.GameObjects.PushBack(obj)
}

func NewGameObj(pos Vec3, size Vec3, color Vec3, mass float32,friction float32,Parent *GameObj)(*GameObj){
	if Parent != nil {
		pos = pos.Add(Parent.Pos)
	}
	//fmt.Println("New game obj friction: " , friction)
	output :=  &GameObj{pos,size,color,Vec3{0,0,0},mass,friction,func(self *Box3D, t float32){},new(list.List),nil,Vec3{0,0,0},*MakeBody(pos,size,mass,friction,false),false }
	if Parent != nil {
		Parent.AddChildren(output)
	} 
	return output
}

type Player struct {
	Body *GameObj
	RightLeg *GameObj
	LeftLeg *GameObj
	RightArm *GameObj
	LeftArm *GameObj
	Head *GameObj

	//constraints
	RightLegConstraint ObjectConstraint
	LeftLegConstraint ObjectConstraint
	RightArmConstraint ObjectConstraint
	LeftArmConstraint ObjectConstraint
	HeadConstraint ObjectConstraint

	//Rest positions
	RightLegRest Vec3
	LeftLegRest Vec3
	RightArmRest Vec3
	LeftArmRest Vec3
	HeadRest Vec3
}

func SetPlayerAnimation (self *Player)(*Player){
	self.Body.Anim = func(box *Box3D, t float32){
		self.HeadConstraint.V = self.HeadRest.Add(Vec3{0,Sin32(t*10),0})
		self.LeftArmConstraint.V = self.LeftArmRest.Add( Vec3{Sin32(t*10)*2,0,0} )
		self.RightArmConstraint.V = self.RightArmRest.Add( Vec3{Sin32(t*10 +float32(math.Pi)/2)*2,0,0} )
		self.LeftLegConstraint.V = self.LeftLegRest.Add( Vec3{0,Sin32(t*10)*4,0} )
		self.RightLegConstraint.V = self.RightLegRest.Add( Vec3{0,Sin32(t*10)*4,0} )
		//self.RightLeg.Body.IsGhost = true
		//self.LeftLeg.Body.IsGhost = true
	}
	return self

} 


func MakePlayer(position Vec3)(newPlayer *Player){
	newPlayer = new(Player)
	newPlayer.Body = NewGameObj(position,Vec3{1,2,1},Vec3{0.5,0.5,0.5},50,0,nil)
	newPlayer.RightLegRest = Vec3{2,-6,0}
	newPlayer.LeftLegRest = Vec3{-2,-6,0}
	newPlayer.RightArmRest = Vec3{2,4,0}
	newPlayer.LeftArmRest = Vec3{-2,4,0}
	newPlayer.HeadRest = Vec3{0,4,0}
	legSize := Vec3{0.5,0.5,0.5}
	newPlayer.RightLeg = NewGameObj(newPlayer.RightLegRest,legSize,Vec3{0.1,0.2,0.3},1.1,5,newPlayer.Body)
	newPlayer.LeftLeg = NewGameObj(newPlayer.LeftLegRest,legSize,Vec3{0.1,0.2,0.3},1,5,newPlayer.Body)
	newPlayer.RightArm = NewGameObj(newPlayer.RightArmRest,legSize,Vec3{0.2,0.1,0.3},1,5,newPlayer.Body)
	newPlayer.LeftArm = NewGameObj(newPlayer.LeftArmRest,legSize,Vec3{0.2,0.1,0.3},1,5,newPlayer.Body)
	newPlayer.Head = NewGameObj(newPlayer.HeadRest,Vec3{1,1,1},Vec3{0.4,0.4,0.4},1,5,newPlayer.Body)

	newPlayer.RightLegConstraint =ObjectConstraint{newPlayer.Body,newPlayer.RightLeg,newPlayer.RightLegRest, 1000,200}
	newPlayer.LeftLegConstraint = ObjectConstraint{newPlayer.Body,newPlayer.LeftLeg,newPlayer.LeftLegRest,1000,200}
	newPlayer.RightArmConstraint = ObjectConstraint{newPlayer.Body,newPlayer.RightArm,newPlayer.RightArmRest,1000,200}

	newPlayer.LeftArmConstraint = ObjectConstraint{newPlayer.Body,newPlayer.LeftArm,newPlayer.LeftArmRest,1000,200}
	newPlayer.HeadConstraint = ObjectConstraint{newPlayer.Body,newPlayer.Head,newPlayer.HeadRest,1000,200}

	newPlayer.Body.Constraints = []Constraint{&(newPlayer.RightLegConstraint),&(newPlayer.LeftLegConstraint),&(newPlayer.RightArmConstraint),&(newPlayer.LeftArmConstraint),&newPlayer.HeadConstraint}
	
	return newPlayer
}



func MakeMan(position Vec3)(*GameObj){
	body := NewGameObj(position,Vec3{1,2,1},Vec3{0.5,0.5,0.5},50,0,nil)

	rleg := NewGameObj(position.Add(Vec3{2,0,0}),Vec3{0.5,0.5,0.5},Vec3{0.1,0.2,0.3},1.1,5,body)
	rlegc := ObjectConstraint{body,rleg,Vec3{2,-6,0},1000,200}
	
	lleg := NewGameObj(Vec3{-2,0,0},Vec3{0.5,0.5,0.5},Vec3{0.3,0.2,0.1},1.1,5,body)
	llegc := ObjectConstraint{body,lleg,Vec3{-2,-6,0},1000,200}
	
	rarm := NewGameObj(Vec3{2,10,0},Vec3{0.5,0.5,0.5},Vec3{0.1,0.2,0.7},1,0,body)
	rarmc := ObjectConstraint{body,rarm,Vec3{2,4,0},100,10}
	
	larm := NewGameObj(Vec3{-2,10,0},Vec3{0.5,0.5,0.5},Vec3{0.3,0.2,0.1},1,0,body)
	larmc := ObjectConstraint{body,larm,Vec3{-2,4,0},100,10}
	
	head := NewGameObj(Vec3{0,5,0},Vec3{1,1,1},Vec3{0.5,0.2,0.2},5,0,body)
	headc := ObjectConstraint{body,head,Vec3{0,4,0},1000,40}
	

	body.Constraints = []Constraint{&rlegc,&llegc,&rarmc, &headc,&larmc}
	body.Anim = func(self *Box3D,t float32)(){
		
	}
	var jumpBusy bool
	jumpBusy = false
	rleg_begin := rlegc.V
	lleg_begin := llegc.V
	/*larm_begin := larmc.V
	rarm_begin := rarmc.V
	var walkCycle float32 = 0
	var speed float32 = 1
	advAnim := func(self *Box3D,t float32){
		walkCycle += 0.01*speed
		
		
		rlegc.V = (rleg_begin.Add(Vec3{0,3,0}.Rotate(walkCycle + math.Pi,0))).Rotate(0, -body.Rotation.X)
		llegc.V = (lleg_begin.Add(Vec3{0,3,0}.Rotate(walkCycle, 0))).Rotate(0, -body.Rotation.X)
		rarmc.V = rarm_begin.Add(Vec3{0,-2,0}.Rotate(walkCycle,0)).Rotate(0,-body.Rotation.X)
		larmc.V = larm_begin.Rotate(0,-body.Rotation.X)
	

}
	body.Anim = advAnim
	*/

	StartJump := func(self *Box3D, t float32){
		jumpStart := t
		if rlegc.V.Y < llegc.V.Y {
			rlegc.V = (rleg_begin.Sub(Vec3{0,5,-5})).Rotate(0,-body.Rotation.X)
		}else {
			llegc.V = (lleg_begin.Sub(Vec3{0,5,-5})).Rotate(0,-body.Rotation.X)
		}
		body.Anim = func(self *Box3D, t float32){
			if t-jumpStart > 0.5 {
				rlegc.V = rleg_begin
				llegc.V = lleg_begin
				
			}
		}
	}
	var direction float32 = 1
	

	Walk := func(self *Box3D, t float32){
		var tc float32 = 0
		
		body.Anim = func(self * Box3D, t float32){
			
			tc = t*10
			rlegc.V = (rleg_begin.Add(Vec3{0,3,0}.Rotate(direction*(tc + math.Pi),0))).Rotate(0, -body.Rotation.X)
			llegc.V = (lleg_begin.Add(Vec3{0,3,0}.Rotate(direction*tc,0))).Rotate(0, -body.Rotation.X)

			
		}

	}

	
	glfw.AddListener(
	func(keyev glfw.KeyEvent){
	fmt.Println(keyev)		
	switch keyev.Key{

		case glfw.KEY_RIGHT, glfw.KEY_D : {
				
			}
		case glfw.KEY_LEFT: {
				
				
			}
		case glfw.KEY_UP: {
				direction = 1
				body.Anim = Walk
				
			}
		case glfw.KEY_DOWN: {
				direction = -1
				body.Anim = Walk
			}
		case glfw.KEY_SPACE: {
				if !jumpBusy{
					body.Anim = StartJump
				}
			}
		}
	})

	return body

}

func testbox(position Vec3)(*GameObj){
	body := NewGameObj(position,Vec3{1,1,1},Vec3{1,0,0},10,0,nil)

	return body
}

func testbox2(position Vec3)(*GameObj){
	
	body := NewGameObj(position,Vec3{1,1,1},Vec3{0,0,1},85,0,nil)
	
	
	return body
}

func ropetest(position Vec3, joints int,dist float32)(*GameObj){
	body := NewGameObj(position,Vec3{1,1,1},Vec3{1,0.5,0.5},Inf,0,nil)
	last := NewGameObj(Vec3{0,-1,0}, Vec3{1,1,1},Vec3{1,0.5,0.5},1.5,1,body)
	lastc := RopeConstraint{body,last,dist*2}
	body.Constraints = []Constraint{lastc}
	
	for i:= 0; i < joints; i++ {
		nlast := NewGameObj(Vec3{0,-dist,0}, Vec3{1,1,1},Vec3{1,0.5,0.5},1.5,10,last)
		nlastc := RopeConstraint{last,nlast,dist*2}
		last.Constraints = []Constraint{nlastc}
		last = nlast
	}

	

	return body

}

func testbox3(position Vec3)(*GameObj){

	body := NewGameObj(position,Vec3{1,1,1},Vec3{0,0,1},84,0,nil)
	//body2 := NewGameObj(Vec3{0,2,0},Vec3{1,1,1},Vec3{1,0,0},1,0,body)
	//	body.Constraints = []Constraint{ObjectConstraint{body,body2,Vec3{0,4,0},100,40}}
	glfw.AddListener(func(keyev glfw.KeyEvent){
		
	})
	return body
}

func treeThing(position Vec3,lvs int)(*GameObj){
	stem := NewGameObj(position,Vec3{1,2,1},Vec3{0.5,0.1,0.3},10,0,nil)
	last := stem
	var i int
	for i= 0; i < lvs; i++ {
		last = NewGameObj(Vec3{0,2,0},Vec3{1,0,1}.Scale(float32(lvs - i)).Add(Vec3{0,1,0}),Vec3{0.2,0.5,0.2},10,0,last)
	}

	return stem
}




func main(){
	//ABSPTest()
	//return
	glfw.Init(800,600)
	InitPhysics()

	world := new(World)
	world.Init()
	world.GameObjects = new(list.List)
	player := MakeMan(Vec3{10,20,10})

	world.Add(player)
	world.Add(ropetest(Vec3{0,40,0},4,4))
	world.Add(treeThing(Vec3{240,20,240},3))
	world.Add(MakePlayer(Vec3{-20,20,0}).Body)
	world.Add(MakePlayer(Vec3{-20,50,0}).Body)
	
	world.Add(MakePlayer(Vec3{-20,70,0}).Body)
	world.Add(MakePlayer(Vec3{-20,90,0}).Body)
	world.Add(MakePlayer(Vec3{-20,110,0}).Body)
	for i := 0; i < 400; i++ {
		world.Add(SetPlayerAnimation(MakePlayer(Vec3{float32(int(i%10))*50,0,float32(i*2) - 200})).Body)
	}
	world.Add(NewGameObj(Vec3{0,-20,0},Vec3{10000,10,10000},Vec3{0,0.5,0.1},float32(math.Inf(1)),10,nil))
	
	//world.Add(SetPlayerAnimation(MakePlayer(Vec3{-20,120,0})).Body)

	qtn := new(ABSPNode)
	qtn.Root = qtn
	//qtn.Position = Vec3{-10000,-10000,-10000}
	//qtn.Size = Vec3{20000,20000,20000}
	for i:= world.GameObjects.Front(); i!= nil; i = i.Next() {
		gobj := i.Value.(*GameObj)
		all := gobj.GetSubs()
		for i:= 0; i < len(all);i++ {
			qtn.Insert(all[i])
		}
	}
	world.GameObjectTree = qtn
	fmt.Println("Total:", len(qtn.Data))
	qtn.Divide()
	cols := 0
	qtn.cd(func(obj1, obj2 SPData){
		cols +=1
	})
	fmt.Println(cols)
	//qtn.Traverse(0)
	//return
	gl.Init()
	vs := gl.CreateShader(gl.VERTEX_SHADER)
	vs.Source(
		LoadFileToString("s1.vert"))
	vs.Compile()

	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	fs.Source( LoadFileToString("s1.frag"))
	fs.Compile()
	
	pg := gl.CreateProgram()
	pg.AttachShader(vs)
	pg.AttachShader(fs)
	pg.Link()
	pg.Validate()
	
	pg.Use()
	fmt.Println("**Shader log**")
	fmt.Println(fs.GetInfoLog())
	fmt.Println(vs.GetInfoLog())
	fmt.Println(pg.GetInfoLog())
	fmt.Println("******END*****")
	
	gl.ClearColor(0.5,0.5,1,0)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.POLYGON_SMOOTH)
	gl.Hint(gl.POLYGON_SMOOTH_HINT,gl.NICEST)

	var t float64
	var ot float64
	var dt float32
	t = float64(time.Nanoseconds())/1000000000
	cam1 := Camera{player,100,Vec3{0,0,0}}
	glfw.AddListener(func(m glfw.MouseMoveEvent){
		cam1.Angle.X = float32(m.X - 400)/400*3.14*2
		cam1.Angle.Y = float32(m.Y - 300)/300*3.14*2
		//player.Rotation = Vec3{cam1.Angle.X,cam1.Angle.Y,0}
	})
	glfw.AddListener(func(mw glfw.MouseWheelEvent){
		cam1.Distance = 100 + float32(mw.Pos*mw.Pos*mw.Pos)
	})


	for it := 0; it < 1000; it +=1 {
		cam1.Setup()
		dt = float32(t - ot)
		
		//fmt.Println(dt)
		ot = t
		dt = 0.01
		t = float64(float64(time.Nanoseconds())/1000000000)
		pt := float64(float64(time.Nanoseconds())/1000000000)
		DoPhysics(world.GameObjects,world.GameObjectTree,dt)//float32(t-ot))
		fmt.Println(float64(float64(time.Nanoseconds())/1000000000) - pt)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		DrawWorld(world.GameObjects,dt,pg)
		glfw.SwapBuffers()
		
	}
	
}
