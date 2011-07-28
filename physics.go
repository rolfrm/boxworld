package main
import "container/list"
import "math"
import "fmt"
func eh(){
	fmt.Println("eh")
}

var Inf float32

func InitPhysics(){
	Inf = float32(math.Inf(1))
}


type Constraint interface{
	Apply(float32)
}


type floatBodySet struct {
	dist float32
	body *PhysicsBody
}

type PhysicsBody struct{
	Pos Vec3
	Size Vec3
	Velocity Vec3
	Mass float32
	Friction float32
	IsGhost bool
	LastVelocity Vec3
}

func MakeBody(Pos Vec3, Size Vec3, Mass float32, Friction float32, IsGhost bool)(nout *PhysicsBody){
	nout = new(PhysicsBody)
	nout.Pos = Pos
	nout.Size = Size
	nout.Mass = Mass
	nout.Friction = Friction
	nout.IsGhost = IsGhost
	return
}


type PhysicsObject interface{
	GetPosition()(Vec3)
	SetPosition(Vec3)
	GetBoundingBox()(Box)
	GetVelocity()(Vec3)
	SetVelocity(Vec3)()
	GetMass()(float32)
	SetMass(float32)
	GetFriction()(float32)
	GetConstraints()([]Constraint)
	GetSubs()([]PhysicsObject)
	GetSize()(Vec3)

	GetBody()(*PhysicsBody)
	UpdatePhysics()

}



type VectorConstraint struct{
	Pos Vec3 
	PObj PhysicsObject
	F func(x Vec3)(Vec3)
}

func (vc VectorConstraint) Apply(dt float32){
	/*diff := vc.Pos.Sub(vc.PObj.GetPosition())
	ApplyImpulse(vc.PObj,vc.F(diff).Scale(dt))
	*/
}



type ObjectConstraint struct{
	P1 PhysicsObject 
	P2 PhysicsObject
	V Vec3
	Spring float32
	Damp float32
	
}

func (oc ObjectConstraint) Apply(dt float32){
	body1 := oc.P1.GetBody()
	body2 := oc.P2.GetBody()
	
	x := body1.Pos.Sub( body2.Pos ).Add(oc.V)
	fspring := x.Scale(oc.Spring)
	dvel :=body1.Velocity.Sub(body2.Velocity)
	fdamp := dvel.Scale(oc.Damp)
	ftotal := Vec3{0,0,0}.Sub(fspring).Sub(fdamp)
	 
	ApplyImpulse(body1,ftotal.Scale(0.5*dt))
	ApplyImpulse(body2,ftotal.Scale(-0.5*dt))

	
}

type RopeConstraint struct{
	P1 PhysicsObject
	P2 PhysicsObject
	Length float32
}
func (oc RopeConstraint) Apply(dt float32){
	var body1 *PhysicsBody = oc.P1.GetBody()
	var body2 *PhysicsBody = oc.P2.GetBody()
	var diff Vec3 = body1.Pos.Sub(body2.Pos)
	var dl float32 = diff.Length()
	
	if dl > oc.Length {
		var dvel Vec3 = body1.Velocity.Sub(body2.Velocity)
		var damp Vec3 = dvel.Scale(100)
		var opf Vec3 = Vec3{0,0,0}.Sub(diff.Scale(20)).Sub(damp)
		ApplyImpulse(body1,opf.Scale(0.5*dt))
		ApplyImpulse(body2,opf.Scale(-0.5*dt))
	}
	 
}


type AABB interface {
	GetBox()(Box)
}

type Box struct {
	Pos Vec3
	Size Vec3
}

func (s Box) GetBox()(Box){
	return s
}
func ApplyImpulse(b1 *PhysicsBody, impulse Vec3){
	b1.Velocity = b1.Velocity.Add(impulse.Scale(1/b1.Mass))
}

func handleCollision(b1 *PhysicsBody, b2 *PhysicsBody, dt float32){
	var diff Vec3 = b1.Pos.Sub(b2.Pos)
	var absDiff Vec3 = diff.Abs()
	var overlap Vec3 = absDiff.Sub(b1.Size.Add(b2.Size))
	
	var collisionAxis int = overlap.BiggestComponent()
	var smallestLength float32
	
	smallestLength = overlap.GetComponent(collisionAxis)

	if smallestLength >= 0 { // Collision?
		return
	} 

	var moveOut Vec3 = Vec3{0,0,0}	
	moveOut.SetComponent(collisionAxis,smallestLength)
	if diff.GetComponent(collisionAxis) < 0 {
		moveOut = moveOut.Scale(-1)
	}
	var n Vec3 = moveOut.Normalize();
	
	totalMass := b1.Mass + b2.Mass
	var b1MoveWeight,b2MoveWeight float32
	if b1.Mass == Inf{
		b2MoveWeight = 1
		b1MoveWeight = 0
	}else if b2.Mass == Inf {
		b1MoveWeight = 1
		b2MoveWeight = 0
	}else{
		b1MoveWeight = b1.Mass/totalMass
		b2MoveWeight = b2.Mass/totalMass
	}
	
	b1.Pos = b1.Pos.Sub(moveOut.Scale(b1MoveWeight))
	b2.Pos = b2.Pos.Add(moveOut.Scale(b2MoveWeight))
	
	difvel := b1.Velocity.Sub(b2.Velocity)

	j :=  difvel.Scale(-1).Dot(n)/(1/b1.Mass + 1/b2.Mass) //Normal force
	if b1.Mass != Inf {
		b1.Velocity = b1.Velocity.Add(n.Scale(j/b1.Mass))
	}
	if b2.Mass != Inf {
		b2.Velocity = b2.Velocity.Sub(n.Scale(j/b2.Mass))
	}
	//Total energy is kept
	var totalFriction float32 = Sqrt32(b1.Friction * b2.Friction)*j
	var tangent Vec3 = Vec3{1,1,1}.Sub(n.Abs())
	var surfSpeed Vec3 = tangent.ElemMul(difvel)
	var damp Vec3 = surfSpeed.Scale(totalFriction).Scale(-dt)
	if damp.Length() > 0 {
		ft :=difvel.ElemMul(tangent).Scale(1/(1/b1.Mass + 1/b2.Mass))		
		if ft.Length() < Fabs32(totalFriction){ //Static
			if b1.Mass != Inf {
				ApplyImpulse(b1,ft.Scale(-b1.Mass))
			}
			if b2.Mass != Inf {
				ApplyImpulse(b2,ft.Scale(b2.Mass))
			}
		}else{ //Dynamic
			ApplyImpulse(b1,damp.Scale(-1))
			ApplyImpulse(b2,damp.Scale(1))
		}   	
	}

}
func createCollisionAtom(dt float32)(func(SPData,SPData)){
	var o1 SPData
	var o1P Vec3
	var o1S Vec3

	var o2 SPData
	var o2P Vec3
	var o2S Vec3
	//Speed bo
 return func (obj1, obj2 SPData){
		if obj1 != o1 { //Memorization optimization
			o1 = obj1
			o1P = obj1.GetPosition()
			o1S = obj1.GetSize()
		}
		if obj2 != o2 {
			o2 = obj2
			o2P = obj2.GetPosition()
			o2S = obj2.GetSize()
		}

		if Fabs32(o2P.X - o1P.X) - (o2S.X + o1S.X) < 0 && Fabs32(o2P.Z - o1P.Z) - (o2S.Z + o1S.Z) < 0 && Fabs32(o2P.Y - o1P.Y) - (o2S.Y  + o1S.Y) < 0 {
		body1,_ := obj1.(*GameObj)
		body2,_ := obj2.(*GameObj)
			if body1 == body2 {
				return
			}
		handleCollision(body1.GetBody(),body2.GetBody(),dt)
	}
		

}
}
func DoPhysics(worldObjects *list.List, objectTree *ABSPNode,dt float32){
	objectTree.Update()
	var allObjects []PhysicsObject
	for item := worldObjects.Front(); item != nil;item = item.Next() {
		ob, ok := item.Value.(PhysicsObject)
		
		if ok{
			allObjects = append(allObjects, ob.GetSubs()...) 
		}
	}
	
	var body1 *PhysicsBody
	var obj1 PhysicsObject
	for i  := 0; i < len(allObjects);i++ {
		obj1 = allObjects[i]
		body1 = obj1.GetBody()
		body1.LastVelocity = body1.Velocity
		if body1.Mass != Inf {
			ApplyImpulse(body1,Vec3{0,-10*dt*body1.Mass,0})
		}
		constraints := obj1.GetConstraints()
		for j:= 0; j < len(constraints); j++ {
			constraints[j].Apply(dt)
		}

		acc := body1.Velocity.Sub(body1.LastVelocity).Scale(dt*dt)
		body1.Pos = body1.Pos.Add(body1.Velocity.Scale(dt)).Add(acc)
		
}
	objectTree.cd(createCollisionAtom(dt))

	for i:= 0; i < len(allObjects);i++ {
		allObjects[i].UpdatePhysics()
	}
}
