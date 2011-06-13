package main

import "gl"
import "container/list"
import "math"
//import "fmt"
type Drawable interface{
	GetBox3D()(Box3D)
	GetChildren()(*list.Element)
}

type Camera struct{
	Target *GameObj
	Distance float32
	Angle Vec3
}

func (self *Camera) Setup(){
	d := self.Distance
	Angle := self.Angle
	Angle360:=Angle.Scale(1/math.Pi*180)
	TargetPos := self.Target.Pos

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Viewport(0,0,800,600)
	gl.Frustum(-1,1 ,-1 ,1 , 4,1000)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(0,0,-d)
	gl.Rotatef(Angle360.Y,1,0,0)
	gl.Rotatef(Angle360.X,0,1,0)
	gl.Translatef(-TargetPos.X,-TargetPos.Y,-TargetPos.Z)
	
}



type Box3D struct{
	Pos Vec3
	Size Vec3
	Color Vec3
	Animator func(*Box3D, float32)
}


var verts []float32 
var verts_init bool
func GetVerts()([]float32){
	if(verts_init == false){
		verts_init = true
		verts = [] float32{1,1,1 ,1,-1,1 ,-1,-1,1, -1,1,1, //Front
		1,1,-1, -1,1,-1, -1,-1,-1, 1,-1,-1, //back

		1,1,1,  -1,1,1, -1,1,-1,  1,1,-1, //top
		1,-1,1,  1,-1,-1, -1,-1,-1,  -1,-1,1, //button
		
		1,1,1  ,1,1,-1 ,1,-1,-1, 1,-1,1, //left
		-1,1,1  ,-1,-1,1 ,-1,-1,-1, -1,1,-1} //righ
	}
	return verts
}



func DrawAABB(aabb Drawable, t float32, program gl.Program){
	box := aabb.GetBox3D()
	nverts := GetVerts()
	if(box.Animator != nil){
		box.Animator(&box,t)
	}
	drawPos :=  box.Pos
	drawColor := box.Color
	drawSize := box.Size
	
	program.GetUniformLocation("SizeVec").Uniform3f(drawSize.X,drawSize.Y,drawSize.Z)
	program.GetUniformLocation("PosVec").Uniform3f(drawPos.X,drawPos.Y,drawPos.Z)
	gl.Color3f(drawColor.X,drawColor.Y,drawColor.Z)
	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.VertexPointer(3,0,nverts)
	gl.DrawArrays(gl.QUADS,0,24)
	gl.EnableClientState(0)
	for it := aabb.GetChildren(); it != nil;it = it.Next() {
		DrawAABB(it.Value.(Drawable),t,program)
	}
}
var draws float32
func DrawWorld(lst *list.List,t float32, program gl.Program){
	draws += t
	for e:= lst.Front(); e != nil ; e = e.Next() {
		val, ok := e.Value.(Drawable)
		if ok {
			DrawAABB(val, draws,program)
		}
	}
}

func DrawBoxes(lst *list.List){


}
