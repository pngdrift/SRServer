attribute vec4 a_position;
attribute vec4 a_color;
attribute vec2 a_texCoord0;

uniform mat4 u_projTrans;
uniform vec4 u_circle;

varying vec4 v_color;
varying vec2 v_texCoords;
varying vec2 v_norm;

void main() {
	v_color = a_color;
	v_color.a = v_color.a * (255.0 / 254.0);
	v_texCoords = a_texCoord0;
	float x = v_texCoords.x - u_circle.x;
	float y = v_texCoords.y - u_circle.y;
	float a = u_circle.z;
	float b = u_circle.w;
	v_norm = vec2(x / a, y / b);
	gl_Position =  u_projTrans * a_position;
}