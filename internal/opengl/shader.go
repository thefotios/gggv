package opengl

import (
	"fmt"

	"github.com/dianelooney/gggv/internal/carbon"
	"github.com/dianelooney/gggv/internal/fps"
	"github.com/go-gl/mathgl/mgl32"
)

const SHADER_TEXTURE_COUNT = 10

type ShaderSource struct {
	name       SourceName
	flipOutput bool

	p        string
	sources  [SHADER_TEXTURE_COUNT]SourceName
	uniforms map[string]BindUniformer

	fbo     uint32
	rbo     uint32
	texture uint32
}

func (s *ShaderSource) Name() SourceName {
	return s.name
}
func (s *ShaderSource) Children() []SourceName {
	out := []SourceName{}
	for i := 0; i < SHADER_TEXTURE_COUNT; i++ {
		if s.sources[i] != "" {
			out = append(out, s.sources[i])
		}
	}
	return out
}
func (s *ShaderSource) Render(scene *Scene) {
	program := scene.programs[s.p].GLProgram
	carbon.BindFramebuffer(carbon.FRAMEBUFFER, s.fbo)
	carbon.UseProgram(program)

	carbon.ActiveTexture(carbon.TEXTURE0)
	carbon.BindTexture(carbon.TEXTURE_2D, s.texture)

	for i, name := range s.sources {
		if name == "" || name == "-" {
			continue
		}

		source, ok := scene.sources[name]
		if !ok {
			continue
		}
		carbon.ActiveTexture(carbon.TEXTURE1 + uint32(i))
		carbon.BindTexture(carbon.TEXTURE_2D, source.Texture())

		switch src := source.(type) {
		case *FFVideoSource:
			carbon.Uniform(program, fmt.Sprintf("tex%vwidth", i), src.width)
			carbon.Uniform(program, fmt.Sprintf("tex%vheight", i), src.height)
		}
	}

	{
		vertAttrib := uint32(carbon.GetAttribLocation(program, carbon.Str("vert\x00")))
		carbon.EnableVertexAttribArray(vertAttrib)
		carbon.VertexAttribPointer(vertAttrib, 3, carbon.FLOAT, false, 5*4, carbon.PtrOffset(0))

		texCoordAttrib := uint32(carbon.GetAttribLocation(program, carbon.Str("vertTexCoord\x00")))
		carbon.EnableVertexAttribArray(texCoordAttrib)
		carbon.VertexAttribPointer(texCoordAttrib, 2, carbon.FLOAT, false, 5*4, carbon.PtrOffset(3*4))

		carbon.Uniform(program, "camera", scene.Camera)

		carbon.UniformTex(program, "lastFrame", 0)
		for i := 0; i < SHADER_TEXTURE_COUNT; i++ {
			carbon.UniformTex(program, fmt.Sprintf("tex%v", i), int32(i)+1)
		}

		carbon.Uniform(program, "time", scene.time)
		carbon.Uniform(program, "fps", float32(fps.LastSec()))
		carbon.Uniform(program, "renderTime", float32(fps.FrameDuration())/NANOSTOSEC)
		x, y := scene.Window.GetCursorPos()
		carbon.Uniform(program, "cursorX", x)
		carbon.Uniform(program, "cursorY", y)
		windowWidth, windowHeight := scene.Window.GetSize()
		carbon.Uniform(program, "windowWidth", windowWidth)
		carbon.Uniform(program, "windowHeight", windowHeight)
		carbon.Uniform(program, "windowSize", [2]float32{float32(windowWidth), float32(windowHeight)})
	}

	for _, u := range scene.uniforms {
		u.BindUniform(program)
	}

	for _, u := range s.uniforms {
		u.BindUniform(program)
	}

	if s.flipOutput {
		carbon.Uniform(program, "flipOutput", int(1))
	} else {
		carbon.Uniform(program, "flipOutput", int(0))
	}

	w, h := s.Dimensions()
	projectionMat := proj(float32(w), float32(h))
	carbon.Uniform(program, "projection", projectionMat)
	r := rect(float32(w), float32(h))
	carbon.BufferData(carbon.ARRAY_BUFFER, len(r)*4, carbon.Ptr(&r[0]), carbon.STATIC_DRAW)
	carbon.DrawArrays(carbon.TRIANGLES, 0, int32(len(r)/5))
	carbon.BindFramebuffer(carbon.FRAMEBUFFER, 0)
}
func (s *ShaderSource) SkipRender(scene *Scene) {}
func (s *ShaderSource) Dimensions() (width, height int32) {
	return 1, 1
}
func (s *ShaderSource) Texture() uint32 {
	return s.texture
}

const sqrt3 = 1.732

func proj(w, h float32) mgl32.Mat4 {
	return mgl32.Ortho(-w/2, w/2, -h/2, h/2, 0.1, 10)
}

const size = 5

func rect(w, h float32) (out []float32) {
	w, h = w/2, h/2
	dx := 2 * w / size
	dy := 2 * h / size
	for x := float32(0); x < size; x++ {
		for y := float32(0); y < size; y++ {
			out = append(out, []float32{
				-w + x*dx, -h + y*dy, 0, 0 + x/size, 0 + y/size,
				-w + (x+1)*dx, -h + y*dy, 0, 0 + (x+1)/size, 0 + y/size,
				-w + x*dx, -h + (y+1)*dy, 0, 0 + x/size, 0 + (y+1)/size,
				//
				-w + (x+1)*dx, -h + y*dy, 0, 0 + (x+1)/size, 0 + y/size,
				-w + x*dx, -h + (y+1)*dy, 0, 0 + x/size, 0 + (y+1)/size,
				-w + (x+1)*dx, -h + (y+1)*dy, 0, 0 + (x+1)/size, 0 + (y+1)/size,
			}...)
		}
	}

	/*
		[]float32{
			-w, -h, 0, 0, 0,
			+w, -h, 0, 1, 0,
			-w, +h, 0, 0, 1,
			//
			+w, -h, 0, 1, 0,
			-w, +h, 0, 0, 1,
			+w, +h, 0, 1, 1,
		}
	*/
	return
}

var hexagon = []float32{
	-1, -1, 0, -1, -1,
	+1, -1, 0, +1, -1,
	+0, +0, 0, +0, +0,
	//
	+1, -1, 0, +1, -1,
	+2, +0, 0, +2, +0,
	+0, +0, 0, +0, +0,
	//
	+2, +0, 0, +2, +0,
	+1, +1, 0, +1, +1,
	+0, +0, 0, +0, +0,
	//
	+1, +1, 0, +1, +1,
	-1, +1, 0, -1, +1,
	+0, +0, 0, +0, +0,
	//
	-1, +1, 0, -1, +1,
	-2, +0, 0, -2, +0,
	+0, +0, 0, +0, +0,
	//
	-2, +0, 0, -2, +0,
	-1, -1, 0, -1, -1,
	+0, +0, 0, +0, +0,
}

func init() {
	for i := 0; i < 18; i++ {
		hexagon[5*i+1] *= sqrt3
		hexagon[5*i+3] = (hexagon[5*i+3] + 1) / 2
		hexagon[5*i+4] = (hexagon[5*i+4] + 1) * sqrt3 / 2
	}
}
