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
	output :=  &GameObj{pos,size,color,Vec3{0,0,0},mass,friction,func(self *Box3D, t float32){},new(list.List),nil,Vec3{0,0,0},*MakeBody(pos,size,mass,friction) }
	if Parent != nil {
		Parent.AddChildren(output)
	} 
	return output
}

func MakeMan(position Vec3)(*GameObj){
	body := NewGameObj(position,Vec3{1,2,1},Vec3{0.5,0.5,0.5},50,0,nil)

	rleg := NewGameObj(position.Add(Vec3{2,0,0}),Vec3{0.5,0.5,0.5},Vec3{0.1,0.2,0.3},1.1,1,body)
	rlegc := ObjectConstraint{body,rleg,Vec3{2,-6,0},1000,200}
	
	lleg := NewGameObj(Vec3{-2,0,0},Vec3{0.5,0.5,0.5},Vec3{0.3,0.2,0.1},1.1,1,body)
	llegc := ObjectConstraint{body,lleg,Vec3{-2,-6,0},1000,200}
	
	rarm := NewGameObj(Vec3{2,10,0},Vec3{0.5,0.5,0.5},Vec3{0.1,0.2,0.7},1,0,body)
	rarmConstraint := ObjectConstraint{body,rarm,Vec3{2,4,0},100,10}
	
	larm := NewGameObj(Vec3{-2,10,0},Vec3{0.5,0.5,0.5},Vec3{0.3,0.2,0.1},1,0,body)
	larmc := ObjectConstraint{body,larm,Vec3{-2,4,0},100,10}
	
	head := NewGameObj(Vec3{0,5,0},Vec3{1,1,1},Vec3{0.5,0.2,0.2},5,0,body)
	headc := ObjectConstraint{body,head,Vec3{0,4,0},1000,40}
	

	body.Constraints = []Constraint{&rlegc,&llegc,&rarmConstraint, &headc,&larmc}
	body.Anim = func(self *Box3D,t float32)(){
		
	}
	DefAnim := func (self *Box3D, t float32){
		/*rlegc.V.X = body.Rotation.Z*2
		rlegc.V.Z = -body.Rotation.X*2
		llegc.V.X = -body.Rotation.Z*2
		llegc.V.Z = body.Rotation.X*2*/


	}
	var jumpBusy bool
	jumpBusy = false
	rleg_begin := rlegc.V
	lleg_begin := llegc.V

	StartJump := func(self *Box3D, t float32){
		jumpStart := t
		if rlegc.V.Y < llegc.V.Y {
			rlegc.V.Y -= 8
			rlegc.V.Z -= 4
		}else {
			llegc.V.Y -=8
			llegc.V.Z -=4
		}
		body.Anim = func(self *Box3D, t float32){
			if t-jumpStart > 0.5 {
				rlegc.V = rleg_begin
				llegc.V = lleg_begin
				
				body.Anim = DefAnim
			}
		}
	}
	var direction float32 = 1
	

	Walk := func(self *Box3D, t float32){
		var tc float32 = 0
		
		body.Anim = func(self * Box3D, t float32){
			
			tc = t*10
			/*fmt.Println(t)
			if(tc > math.Pi*2){
				tc = 0
			}*/
			rlegc.V = (rleg_begin.Add(Vec3{0,3,0}.Rotate(direction*(tc + math.Pi),0))).Rotate(0, -body.Rotation.X)
			llegc.V = (lleg_begin.Add(Vec3{0,3,0}.Rotate(direction*tc,0))).Rotate(0, -body.Rotation.X)
			//rlegc.V = (rleg_begin.Add(Vec3{0,3,0}.Rotate(direction*(tc + math.Pi),0)))//.Rotate(0, -body.Rotation.X)
			//llegc.V = (lleg_begin.Add(Vec3{0,3,0}.Rotate(direction*tc,0)))//.Rotate(0, -body.Rotation.X)
			
		}

	}

	
	glfw.AddListener(
	func(keyev glfw.KeyEvent){
	fmt.Println(keyev)		
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
				direction = 1
				body.Anim = Walk
				
			}
		case glfw.KEY_DOWN: {
				direction = -1
				body.Anim = Walk
				ApplyImpulse2(body.GetBody(),body.Rotation.Scale(-body.Mass))
			}
		case glfw.KEY_SPACE: {
				fmt.Println("Eh?")
				if !jumpBusy{
					body.Anim = StartJump
				}
			}
		}
	})
	 
	//body.AddChildren(rleg,lleg, rarm,larm,head)

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
		ApplyImpulse(body,Vec3{0,10*body.Mass,0})
	})
	return body
}


func main(){
	//BSPTest()
	//return
	//CamFocus Vec3

	glfw.Init(800,600)
	InitPhysics()

	ground := new(GameObj)
	ground.Pos = Vec3{0,-3,0}
	ground.Size = Vec3{1000,10,1000}
	ground.Color = Vec3{0,0.5,0.1}
	ground.Mass = float32(math.Inf(1))
	ground.Children = new(list.List)
	ground.Friction = 1
	ground = NewGameObj(Vec3{0,-3,0},Vec3{1000,10,1000},Vec3{0,0.5,0.1},float32(math.Inf(1)),0.1,nil)
	world := new(World)
	world.Init()
	world.GameObjects = new(list.List)
	//world.Add(ground)
	//world.Add(testbox(Vec3{10,0,0}))
	//world.Add(testbox(Vec3{10,10,0}))
	//world.Add(testbox(Vec3{10,15,0}))
	//supertestbox := testbox3(Vec3{-10,10,0})
	//world.Add(supertestbox)
	player := MakeMan(Vec3{10,20,10})
	world.Add(player)
	world.Add(NewGameObj(Vec3{0,-20,0},Vec3{10000,10,10000},Vec3{0,0.5,0.1},float32(math.Inf(1)),10,nil))
	//world.Add(NewGameObj(Vec3{0,150,0},Vec3{10000,10,10000},Vec3{0.5,0.5,0.9},float32(math.Inf(1)),0,nil))
	
	world.Add(ropetest(Vec3{0,40,0},10,4))
	//nbox := testbox2(Vec3{15,20,0})
	//nbox.Mass = 10
	//world.Add(nbox)
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
	
	//gl.ClearColor(0.5,0.5,1,0)
	gl.ClearColor(0,0,0,0)
	//gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.FOG)
	//gl.Disable(gl.DEPTH_TEST)
	//gl.CullFace(gl.BACK)
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
		player.Rotation = Vec3{cam1.Angle.X,cam1.Angle.Y,0}
	})
	glfw.AddListener(func(mw glfw.MouseWheelEvent){
		cam1.Distance = 100 + float32(mw.Pos*mw.Pos*mw.Pos)
	})


	for it := 0; it < 100000; it +=1 {
		cam1.Setup()
		dt = float32(t - ot)
		ot = t
		math.Sin(float64(dt))
		t = float64(float64(time.Nanoseconds())/1000000000)
		DoPhysics(world.GameObjects,0.001)//float32(t-ot))
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		DrawWorld(world.GameObjects,0.001,pg)
		glfw.SwapBuffers()
		time.Sleep(100000)
	}
	
}
