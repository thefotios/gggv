#version 330

//default uniforms from GGGV

//futher sources are set as tex1, tex2, tex3
uniform sampler2D tex0;

//result of the last render of this shader
uniform sampler2D lastFrame;

//only set when tex* is a FFsource
uniform float tex0Width;
uniform float tex0Height;

//seconds since application start. 
//monotonically increasing, will never be 0 
uniform float time;
//time the last frame took to render
uniform float renderTime;
//number of frames in the last second
uniform float fps;
//output window size
uniform float windowHeight;
uniform float windowWidth;
//cursor position
uniform float cursorX;
uniform float cursorY;

//texCoord
in vec2 fragTexCoord;
//output pixel color
out vec4 outputColor;

uniform float fadeRate = 0.003;
uniform float threshold = 0.4;

uniform float highlightR = 255;
uniform float highlightG = 255;
uniform float highlightB = 255;

void main() {
  outputColor = texture(tex0, fragTexCoord);
  float total = (outputColor.r + outputColor.g + outputColor.b) / 3;

  if (total > threshold) {
    outputColor = vec4(1, 1, 1, 1);
  } else {
    outputColor = vec4(0, 0, 0, 1);
  }
}
