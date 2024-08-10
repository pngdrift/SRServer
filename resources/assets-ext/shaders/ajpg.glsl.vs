attribute vec4 a_position;
attribute vec4 a_color;
attribute vec2 a_texCoord0;

uniform mat4 u_projTrans;

varying vec4 v_color;
varying vec2 v_texCoordsColor0;
varying vec2 v_texCoordsAlpha0;

void main() {
	v_color = a_color;
	v_color.a = v_color.a * (255.0 / 254.0);
	v_texCoordsColor0 = vec2(a_texCoord0.x, a_texCoord0.y * 0.5);
	v_texCoordsAlpha0 = vec2(a_texCoord0.x, a_texCoord0.y * 0.5 + 0.5);
	gl_Position =  u_projTrans * a_position;
}