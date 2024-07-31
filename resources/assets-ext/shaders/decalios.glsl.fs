#ifdef GL_ES
#define LOWP lowp
precision highp float;
#else
#define LOWP
#endif

varying LOWP vec4 v_color;
varying vec2 v_texCoords0;
varying vec2 v_texCoords1;

uniform sampler2D u_texture0;
uniform sampler2D u_texture1;

uniform float u_H;
uniform float u_S;
uniform float u_L;



void main() {
	/* CAR */
	vec4 tc0 = texture2D(u_texture0, v_texCoords0);
	/* DECAL */
	vec4 tc1 = v_color * texture2D(u_texture1, v_texCoords1);
	vec4 color = vec4(tc1.r, tc1.g, tc1.b, tc1.a);

	float alpha = min(tc0.a, tc1.a);
	color.a = alpha;
	gl_FragColor = color;
}