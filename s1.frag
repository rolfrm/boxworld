//Fragment shader..
varying float fogFactor;
varying vec3 fragPos;
void main(){
  float fragPosition = length(fragPos);
  //fragPosition = 0;
  //gl_Color.r = sin(fragPos.x)*0.5 + 0.5;
  gl_FragColor = mix(vec4(0.5,0.5,1,1), gl_Color,exp(-fragPosition/300));
}
