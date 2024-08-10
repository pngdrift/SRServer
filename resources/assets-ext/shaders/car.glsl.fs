#ifdef GL_ES
#define LOWP lowp
precision highp float;
#else
#define LOWP
#endif

varying LOWP vec4 v_color;
varying vec2 v_texCoords;

uniform sampler2D u_texture;
uniform vec4 u_car_color;

void main() {
	vec4 tc = v_color * texture2D(u_texture, v_texCoords);
	float gray = 0.21 * tc.r + 0.72 * tc.g + 0.07 * tc.b;
	float alpha = tc.a;
	vec4 color = vec4(u_car_color.r, u_car_color.g, u_car_color.b, alpha);
	float shadow = (alpha * pow(1.0 - gray, 4.0));
	color = mix(color, vec4(0.0, 0.0, 0.0, 1.0), shadow);
	float light = (alpha * pow(gray, 2.0));
	color = mix(color, vec4(1.0, 1.0, 1.0, 1.0), light);
	color.a = alpha * u_car_color.a;
	gl_FragColor = color;
}
			