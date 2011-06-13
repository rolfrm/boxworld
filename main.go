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
	Anim func(self *Box3D,t float32)
	Children *list.List
	Constraints []Constraint
	Rotation Vec3

}

func (self *GameObj) GetPosition()(Vec3){
	return self.Pos
}

func (self * GameObj) SetPosition(in Vec3){
	self.Pos = in
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


func (s *GameObj) KeyDown(key int) {
	if key == 'S' {
		s.Pos.X -=1
	}
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
	//fmt.Println(self.GetMass())
	return output
}




type World struct{
	GameObjects *list.List
}

func (self *World) Init(){
	fmt.Println("Hey")
	self.GameObjects = new(list.List)
}

func (self *World) Add(obj * GameObj){
	self.GameObjects.PushBack(obj)
}

func NewGameObj(pos Vec3, size Vec3, color Vec3, mass float32)(GameObj){
	return GameObj{pos,size,color,Vec3{0,0,0},mass,func(self *Box3D, t float32){},new(list.List),nil,Vec3{0,0,0} }
}

func MakeMan(position Vec3)(*GameObj){
	body := new(GameObj)
	*body = NewGameObj(position,Vec3{1,2,1},Vec3{0.5,0.5,0.5},50)

	rleg := NewGameObj(Vec3{2,0,0},Vec3{0.5,0.5,0.5},Vec3{0.1,0.2,0.3},10)
	rlegc := ObjectConstraint{body,&rleg,Vec3{2,-4,0},func(x Vec3)(Vec3){return x.Scale(1000)}}
	lleg := NewGameObj(Vec3{-2,0,0},Vec3{0.5,0.5,0.5},Vec3{0.3,0.2,0.1},10)
	llegc := ObjectConstraint{body,&lleg,Vec3{-2,-4,0},func(x Vec3)(Vec3){return x.Scale(1000)}}
	rarm := NewGameObj(Vec3{2,10,0},Vec3{0.5,0.5,0.5},Vec3{0.1,0.2,0.7},5)
	rarmConstraint := ObjectConstraint{body,&rarm,Vec3{2,4,0},func(x Vec3)(Vec3){
			return x.Scale(1000)
}}
	larm := NewGameObj(Vec3{-2,10,0},Vec3{0.5,0.5,0.5},Vec3{0.3,0.2,0.1},1)
	head := NewGameObj(Vec3{0,10,0},Vec3{1,1,1},Vec3{0.5,0.2,0.2},5)
	
	body.Constraints = []Constraint{&rlegc,&llegc,
		RopeConstraint{body,&larm,2},
	&rarmConstraint,
		ObjectConstraint{body,&head,Vec3{0,4,0},func(x Vec3)(Vec3){return x.Scale(1000)}}}
	body.Anim = func(self *Box3D,t float32)(){
		rarmConstraint.V.X = 2 + float32(math.Sin(float64(body.Pos.X)))
		//fmt.Println(body.Rotation)
	}
	DefAnim := func (self *Box3D, t float32){
		rlegc.V.X = body.Rotation.Z*2
		rlegc.V.Z = -body.Rotation.X*2
		llegc.V.X = -body.Rotation.Z*2
		llegc.V.Z = body.Rotation.X*2


	}
	var jumpBusy bool

	StartJump := func(self *Box3D, t float32){
		jumpStart := t
		jumpBusy = true
		begPoint := rlegc.V
		rlegc.V.Y -= 8
		body.Anim = func(self *Box3D, t float32){
			if t-jumpStart > 0.5 {
				rlegc.V = begPoint
				body.Anim = DefAnim
				jumpBusy = false
			}
		}
	}
	
	glfw.AddListener(
	func(keyev glfw.KeyEvent){
		switch keyev.Key{

		case glfw.KEY_RIGHT, glfw.KEY_D : {
				if keyev.Action == 1 {
					ApplyImpulse(body,Vec3{1*body.Mass,0,0})
				}
			}
		case glfw.KEY_LEFT: {
				ApplyImpulse(body,Vec3{-1*body.Mass,0,0})
			}
		case glfw.KEY_UP: {
				ApplyImpulse(body,body.Rotation.Scale(body.Mass))
			}
		case glfw.KEY_DOWN: {
				ApplyImpulse(body,body.Rotation.Scale(-body.Mass))
			}
		case glfw.KEY_SPACE: {
				if !jumpBusy{
					body.Anim = StartJump
				}
			}
		}
	})
	 
	body.AddChildren(&rleg,&lleg, &rarm,&larm,&head)

	return body

}

func testbox(position Vec3)(*GameObj){
	body := new(GameObj)
	*body = NewGameObj(position,Vec3{1,1,1},Vec3{1,0,0},10)

	return body
}

func testbox2(position Vec3)(*GameObj){
	body := new(GameObj)
	*body = NewGameObj(position,Vec3{1,1,1},Vec3{0,0,1},85)
	
	
	return body
}

func ropetest(position Vec3, joints int,dist float32)(*GameObj){
	body:= new(GameObj)
	*body = NewGameObj(position,Vec3{1,1,1},Vec3{1,0.5,0.5},Inf)
	last := new(GameObj)
	*last = NewGameObj(position.Add(Vec3{0,-dist,0}), Vec3{1,1,1},Vec3{1,0.5,0.5},1)
	lastc := RopeConstraint{body,last,dist*2}
	body.Constraints = []Constraint{lastc}
	body.AddChildren(last)

	for i:= 0; i < joints; i++ {
		nlast := new(GameObj)
		*nlast = NewGameObj(last.Pos.Add(Vec3{0,-dist,0}), Vec3{1,1,1},Vec3{1,0.5,0.5},1)
		nlastc := RopeConstraint{last,nlast,dist*2}
		last.Constraints = []Constraint{nlastc}
		last.AddChildren(nlast)
		last = nlast
	}

	

	return body

}



func testbox3(position Vec3)(*GameObj){
	body := new(GameObj)
	*body = NewGameObj(position,Vec3{1,1,1},Vec3{0,0,1},84)
	body2 := NewGameObj(Vec3{position.X,position.Y+2,position.Z},Vec3{1,1,1},Vec3{1,0,0},1)
	body.Constraints = []Constraint{ObjectConstraint{body,&body2,Vec3{0,4,0},func(x Vec3)(Vec3){
				return x.Scale(1000)
	}}}
	glfw.AddListener(func(keyev glfw.KeyEvent){
		pwr := float32(10*body.GetMass())
		switch keyev.Key {
		case 32: {ApplyImpulse(body,Vec3{0,pwr,0})}
		}
	})

	body.AddChildren(&body2)
	return body
}


func main(){

	//CamFocus Vec3

	glfw.Init(800,600)
	InitPhysics()

	ground := new(GameObj)
	ground.Pos = Vec3{0,-3,0}
	ground.Size = Vec3{100,10,100}
	ground.Color = Vec3{0,0.5,0.1}
	ground.Mass = float32(math.Inf(1))
	ground.Children = new(list.List)
	world := new(World)
	world.Init()
	world.GameObjects = new(list.List)
	world.Add(ground)
	world.Add(testbox(Vec3{10,0,0}))
	world.Add(testbox(Vec3{10,10,0}))
	world.Add(testbox(Vec3{10,15,0}))
	//world.Add(testbox2(Vec3{10,20,0}))
	world.Add(testbox3(Vec3{-10,10,0}))
	player := MakeMan(Vec3{0,10,0})
	world.Add(player)
	world.Add(ropetest(Vec3{0,40,0},10,1.5))
	nbox := testbox2(Vec3{15,20,0})
	nbox.Mass = 10
	world.Add(nbox)
	//world.Add(MakeMan(Vec3{-4,3,0}))
	math.Sin(1)
	glfw.Version()
	
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
	fmt.Println(fs.GetInfoLog())
	fmt.Println(vs.GetInfoLog())
	fmt.Println(pg.GetInfoLog())
	
	gl.ClearColor(0.5,0.5,1,0)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	var t float64
	var ot float64
	t = float64(time.Nanoseconds())/1000000000
	cam1 := Camera{player,100,Vec3{0,0,0}}
	glfw.AddListener(func(m glfw.MouseMoveEvent){
		cam1.Angle.X = float32(m.X - 400)/400*3.14
		cam1.Angle.Y = float32(m.Y - 300)/300*3.14
		player.Rotation = Vec3{float32(math.Sin(float64(cam1.Angle.X))),0,-float32(math.Cos(float64(cam1.Angle.X)))}
	})
	glfw.AddListener(func(mw glfw.MouseWheelEvent){
		fmt.Println(mw.Pos)
		cam1.Distance = 100 + float32(mw.Pos*mw.Pos*mw.Pos)
	})


	for it := 0; it < 100000; it +=1 {
		cam1.Setup()
		ot = t
		t = float64(float64(time.Nanoseconds())/1000000000)
		//fmt.Println(t - ot)
		go DoPhysics(world.GameObjects,float32(t-ot))
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		DrawWorld(world.GameObjects,float32(t-ot),pg)
		glfw.SwapBuffers()
		time.Sleep(1000000)
	}
	
}