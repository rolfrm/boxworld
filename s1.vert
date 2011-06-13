uniform vec3 SizeVec;
uniform vec3 PosVec;

void main(){
  vec4 pos = gl_Vertex;
  pos.x *=SizeVec.x;
  pos.y *=SizeVec.y;
  pos.z *=SizeVec.z;
  pos.xyz += PosVec;
  gl_FrontColor = gl_Color;
  gl_Position = gl_ModelViewProjectionMatrix*pos;
}
