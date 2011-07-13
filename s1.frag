//Fragment shader..
varying float fogFactor;
varying vec3 fragPos;
void main(){
  float fragPosition = length(fragPos);
  gl_FragColor = mix(vec4(0.5,0.5,1,1),gl_Color,exp(-fragPosition/300));
}
