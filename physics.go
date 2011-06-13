package main
import "container/list"
import "math"
var Inf float32

func InitPhysics(){
	Inf = float32(math.Inf(1))
}


type Constraint interface{
	Apply(float32)
}


type PhysicsObject interface{
	GetPosition()(Vec3)
	SetPosition(Vec3)
	GetBoundingBox()(Box)
	GetVelocity()(Vec3)
	SetVelocity(Vec3)()
	GetMass()(float32)
	SetMass(float32)
	GetConstraints()([]Constraint)
	GetSubs()([]PhysicsObject)
	GetSize()(Vec3)
}



type VectorConstraint struct{
	Pos Vec3 
	PObj PhysicsObject
	F func(x Vec3)(Vec3)
}

func (vc VectorConstraint) Apply(dt float32){
	diff := vc.Pos.Sub(vc.PObj.GetPosition())
	ApplyImpulse(vc.PObj,vc.F(diff).Scale(dt))
	
}



type ObjectConstraint struct{
	P1 PhysicsObject 
	P2 PhysicsObject
	V Vec3
	F func(x Vec3 )(Vec3)
}

func (oc ObjectConstraint) Apply(dt float32){
	diff := oc.P2.GetPosition().Sub(oc.P1.GetPosition().Add(oc.V))
	dl := oc.F(diff).Scale(dt)
	dvel :=oc.P2.GetVelocity().Sub(oc.P1.GetVelocity())
	ftotal := dl.Sub(dvel.Scale(-0.9))
	ApplyImpulse(oc.P1,ftotal.Scale(0.5))
	ApplyImpulse(oc.P2,ftotal.Scale(-0.5))
}

type RopeConstraint struct{
	P1 PhysicsObject
	P2 PhysicsObject
	Length float32
}
func (oc RopeConstraint) Apply(dt float32){
	diff := oc.P2.GetPosition().Sub(oc.P1.GetPosition())
	dl := diff.Length()
	
	if dl > oc.Length {
		dvel :=oc.P2.GetVelocity().Sub(oc.P1.GetVelocity())
		damp := dvel.Scale(-0.9)
		opf := diff.Sub(damp).Scale(0.9)
		ApplyImpulse(oc.P1,opf.Scale(0.5))
		ApplyImpulse(oc.P2,opf.Scale(-0.5))
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
func ApplyImpulse(pobj PhysicsObject,impulse Vec3){
	vel := pobj.GetVelocity() //A reference
	vel = vel.Add(impulse.Scale(float32(1)/(pobj.GetMass())))
	pobj.SetVelocity(vel)
}

func CheckCollision(aabb1 AABB, aabb2 AABB) (bool,Vec3){
	b1 := aabb1.GetBox()
	b2 := aabb2.GetBox()
	diff := b2.Pos.Sub(b1.Pos)
	colVec :=diff.ElemDiv(b1.Size.Add(b2.Size))
	collisionCheck := colVec.Abs()
	return  (collisionCheck.X < 1) && (collisionCheck.Y < 1) && (collisionCheck.Z < 1) , colVec
	
}


func CheckCollision2(p1 PhysicsObject, p2 PhysicsObject){

	aabb1 := p1.GetBoundingBox()
	aabb2 := p2.GetBoundingBox()
	col, overlap := CheckCollision(aabb1,aabb2)
	if !col {
		return
	}
	difval := p1.GetVelocity().Sub(p2.GetVelocity())
	if difval.Length() <= 0 {
		return
	}

	var n Vec3
	switch overlap.Abs().BiggestComponent() {
	case 0: n = Vec3{1,0,0} 
	case 1: n = Vec3{0,1,0}
	case 2: n = Vec3{0,0,1}
	}

	nd := p1.GetPosition().ElemMul(n).Sub(p2.GetPosition().ElemMul(n)) //Get Overlap on collision axis
	if nd.Dot(n) < 0 {
		n = n.Scale(-1)
	}
	move := nd.Sub(n.ElemMul(p1.GetSize())).Sub(n.ElemMul(p2.GetSize()))
	m1 := move.Scale(-p2.GetMass()/(p2.GetMass() + p1.GetMass()))
	m2 := move.Scale(p1.GetMass()/(p2.GetMass() + p1.GetMass()))
	if p2.GetMass() == Inf {
		m1 = move.Scale(-1)
	}else if p1.GetMass() == Inf {
		m2 = move.Scale(1)
	}

	p1.SetPosition(p1.GetPosition().Add(m1))
	p2.SetPosition(p2.GetPosition().Add(m2))




	e := float32(0)
	j :=  difval.Scale(-1*(1 + e)).Dot(n)/(1/p1.GetMass() + 1/p2.GetMass())
	p1.SetVelocity(p1.GetVelocity().Add(n.Scale(j/p1.GetMass() )))
	p2.SetVelocity(p2.GetVelocity().Sub(n.Scale(j/p2.GetMass() )))
	

}

func DoPhysics(worldObjects *list.List,dt float32){
	allObjects := []PhysicsObject{}

	for item := worldObjects.Front(); item != nil;item = item.Next() {
		ob, ok := item.Value.(PhysicsObject)
		if ok{
			allObjects = append(allObjects, ob.GetSubs()...) 
		}
	}
	for i := 0; i < len(allObjects);i++ {
		obj1 := allObjects[i]
		
		if obj1.GetMass() != Inf {
			ApplyImpulse(obj1,Vec3{0,-10*dt*obj1.GetMass(),0})
		}
		constraints := obj1.GetConstraints()
		for j:= 0; j < len(constraints); j++ {
			constraints[j].Apply(dt)
		}
			

		obj1.SetPosition(obj1.GetPosition().Add(obj1.GetVelocity().Scale(dt)))
	
		for j := 0; j < len(allObjects);j++ {
			if j == i {
				continue
			}

			CheckCollision2(obj1,allObjects[j])
		}

	}
	
}
