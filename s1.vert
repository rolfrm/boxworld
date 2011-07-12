uniform vec3 SizeVec;
uniform vec3 PosVec;
varying float fogFactor;
varying vec3 fragPos;
void main(){
  
  vec3 norm = gl_Normal;
  vec4 pos = gl_Vertex;
  pos.x *=SizeVec.x;
  pos.y *=SizeVec.y;
  pos.z *=SizeVec.z;
  pos.xyz += PosVec;
  vec4 vVertex = gl_ModelViewMatrix*pos;
  fragPos = vVertex.xyz;
  gl_FogFragCoord = length(vVertex.xyz);
  float fogDensity = 0.1;
  fogFactor = exp(-gl_FogFragCoord/100);//exp2(-fogDensity*gl_FogFragCoord*gl_FogFragCoord*1.442695);
  gl_FrontColor = gl_Color*(0.5*dot(norm,vec3(1,1,1))+ 0.8);
  //gl_FrontColor = gl_Color;
  gl_Position = gl_ProjectionMatrix*vVertex;
}
