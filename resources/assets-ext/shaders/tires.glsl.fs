#ifdef GL_ES
#define LOWP lowp
precision highp float;
#else
#define LOWP
#endif

varying vec4 v_color;
varying vec2 v_texCoords;
varying vec2 v_norm;

uniform sampler2D u_texture;

void main() {
	float nx = v_norm.x;
	float ny = v_norm.y;
	vec4 tc = v_color * texture2D(u_texture, v_texCoords);

    float dst = length(v_norm);
    float r = (nx * nx) + (ny * ny);
    tc.a *= 1.0 - smoothstep(r, r + 0.2, dst);

    gl_FragColor = tc;
//	float r = (nx * nx) + (ny * ny);
//	if (r < 1.0) {
//		tc.a = 0.0;
//	}
//	gl_FragColor = vec4(tc.r, tc.g, tc.b, tc.a);
}
