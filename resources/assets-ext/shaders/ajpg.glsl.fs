#ifdef GL_ES
#define LOWP lowp
precision highp float;
#else
#define LOWP
#endif

varying LOWP vec4 v_color;
varying vec2 v_texCoordsColor0;
varying vec2 v_texCoordsAlpha0;

uniform sampler2D u_texture;

void main() {
	vec4 color = texture2D(u_texture, v_texCoordsColor0);
    vec4 alpha = texture2D(u_texture, v_texCoordsAlpha0);
    vec4 fragment = v_color * vec4(color.r, color.g, color.b, (alpha.r + alpha.g + alpha.b) / 3.0);
	gl_FragColor = fragment;
}