//Fragment shader..
varying float fogFactor;
varying vec3 fragPos;
void main(){
  float fragPosition = length(fragPos);
  gl_FragColor = mix(vec4(0,0,0,1),gl_Color,1);//exp(-fragPosition/100));
}
