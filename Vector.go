package main
import "math"

type Vec3 struct{
	X float32
	Y float32
	Z float32
}

func (self Vec3) Sub(with Vec3)(output Vec3){
	output.X = self.X - with.X
	output.Y = self.Y - with.Y
	output.Z = self.Z - with.Z
	return
}

func (self Vec3) Add(b Vec3)(output Vec3){
	output.X = b.X + self.X
	output.Y = b.Y + self.Y
	output.Z = b.Z + self.Z
	return
}

func (self Vec3) ElemDiv(b Vec3)(output Vec3){
	output.X = self.X/b.X
	output.Y = self.Y/b.Y
	output.Z = self.Z/b.Z
	return
}

func (self Vec3) ElemMul(v Vec3)(Vec3){
	return Vec3{self.X*v.X, self.Y*v.Y, self.Z*v.Z}
}


func Fabs32(in float32)(float32){
	if in < 0 {
		return -in
	}
	return in	
}

func (self Vec3) Abs()(Vec3){
		return Vec3{Fabs32(self.X),Fabs32(self.Y),Fabs32(self.Z)}
}
func (self Vec3) Rotate(rx float32, ry float32) Vec3{
	cosa := Cos32(rx)
	sina := Sin32(rx)
	cosb := Cos32(ry)
	sinb := Sin32(ry)
	x := cosb*self.X + sinb*self.Z
	y := sina*sinb*self.X + self.Y*cosa -self.Z*sina*cosb
	z := -self.X*cosa*sinb + self.Y*sina + self.Z*cosa*cosb
	return Vec3{x,y,z}
}


func signOfFloat32(in float32)(float32){
	if(in < 0){
		return -1
	}
	return 1
}

func (self Vec3) Sign()(Vec3){
	return Vec3{signOfFloat32(self.X),signOfFloat32(self.Y),signOfFloat32(self.Z)}
}

func (self Vec3) NegSign()(Vec3){
	return Vec3{signOfFloat32(-self.X),signOfFloat32(-self.Y),signOfFloat32(-self.Z)}
}

func (self Vec3) Scale( scalator float32)(Vec3){
	return Vec3{self.X*scalator,self.Y*scalator,self.Z*scalator}
}

func (self Vec3) Dot(other Vec3)(float32){
	return self.X*other.X + self.Y*other.Y + self.Z*other.Z
}

func (self Vec3) Length()(float32){
	return float32(math.Sqrt(float64(self.X*self.X + self.Y*self.Y + self.Z*self.Z)))
}

func (self Vec3) Normalize()(Vec3){
	len := self.Length()
	return Vec3{self.X/len, self.Y/len,self.Z/len}
}

func (s Vec3) SmallestComponent()(int){
	if s.X < s.Y && s.X < s.Z {
		return 0
	} else if s.Y < s.X && s.Y < s.Z {
		return 1
	}
	return 2
}
func (s Vec3) BiggestComponent()(int){
	if s.X > s.Y && s.X > s.Z {
		return 0
	} else if s.Y > s.X && s.Y > s.Z {
		return 1
	}
	return 2
}

