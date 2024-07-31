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

float min3(float a, float b, float c) {
    return min(min(a, b), c);
}

float max3(float a, float b, float c) {
    return max(max(a, b), c);
}

vec4 convertRGBtoHSL( vec4 col )
{
    float red   = col.r;
    float green = col.g;
    float blue  = col.b;

    float minc  = min3( col.r, col.g, col.b );
    float maxc  = max3( col.r, col.g, col.b );
    float delta = maxc - minc;

    float lum = (minc + maxc) * 0.5;
    float sat = 0.0;
    float hue = 0.0;

    if (lum > 0.0 && lum < 1.0) {
        float mul = (lum < 0.5)  ?  (lum)  :  (1.0-lum);
        sat = delta / (mul * 2.0);
    }

    vec3 masks = vec3(
        (maxc == red   && maxc != green) ? 1.0 : 0.0,
        (maxc == green && maxc != blue)  ? 1.0 : 0.0,
        (maxc == blue  && maxc != red)   ? 1.0 : 0.0
    );

    vec3 adds = vec3(
              ((green - blue ) / delta),
        2.0 + ((blue  - red  ) / delta),
        4.0 + ((red   - green) / delta)
    );

    float deltaGtz = (delta > 0.0) ? 1.0 : 0.0;

    hue += dot( adds, masks );
    hue *= deltaGtz;
    hue /= 6.0;

    if (hue < 0.0)
        hue += 1.0;

    return vec4( hue, sat, lum, col.a );
}

vec4 convertHSLtoRGB( vec4 col )
{
    const float onethird = 1.0 / 3.0;
    const float twothird = 2.0 / 3.0;
    const float rcpsixth = 6.0;

    float hue = col.x;
    float sat = col.y;
    float lum = col.z;

    vec3 xt = vec3(
        rcpsixth * (hue - twothird),
        0.0,
        rcpsixth * (1.0 - hue)
    );

    if (hue < twothird) {
        xt.r = 0.0;
        xt.g = rcpsixth * (twothird - hue);
        xt.b = rcpsixth * (hue      - onethird);
    }

    if (hue < onethird) {
        xt.r = rcpsixth * (onethird - hue);
        xt.g = rcpsixth * hue;
        xt.b = 0.0;
    }

    xt = min( xt, 1.0 );

    float sat2   =  2.0 * sat;
    float satinv =  1.0 - sat;
    float luminv =  1.0 - lum;
    float lum2m1 = (2.0 * lum) - 1.0;
    vec3  ct     = (sat2 * xt) + satinv;

    vec3 rgb;
    if (lum >= 0.5)
         rgb = (luminv * ct) + lum2m1;
    else rgb =  lum    * ct;

    return vec4( rgb, col.a );
}

void main() {
	/* CAR */
	vec4 tc0 = texture2D(u_texture0, v_texCoords0);
	/* DECAL */
	vec4 tc1 = v_color * texture2D(u_texture1, v_texCoords1);
	vec4 color = vec4(tc1.r, tc1.g, tc1.b, tc1.a);

    // -> to HSL
	vec4 hslColor = convertRGBtoHSL(color);
	/* H ajust */
	hslColor.r += u_H;
	/* S ajust */
	hslColor.g *= u_S;
	/* L ajust */
	hslColor.b *= u_L;

	// -> to RGB
	color = convertHSLtoRGB(hslColor);

	float alpha = min(tc0.a, tc1.a);
	color.a = alpha;
	gl_FragColor = color;
}