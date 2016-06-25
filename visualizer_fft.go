package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/cmplx"
	"os/exec"

	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
)

const FFT_PALETTE_SIZE = 768

type FFTVisualizer struct {
	BaseVisualizer

	Luminosity float64 `property:"rw"`

	palette []Led
	cmd     *exec.Cmd

	samplesPerFrame int
	samplesPerFFT   int

	pulseAudioDevice string
	pulseAudioOpts   []string
}

func NewFFTVisualizer(pulseAudioDevice string, pulseAudioOpts []string) *FFTVisualizer {
	v := FFTVisualizer{
		BaseVisualizer:   *NewBaseVisualizer("FFT"),
		palette:          make([]Led, FFT_PALETTE_SIZE),
		Luminosity:       1.0,
		samplesPerFrame:  512,
		samplesPerFFT:    4096,
		pulseAudioDevice: pulseAudioDevice,
		pulseAudioOpts:   pulseAudioOpts,
	}

	// Generate first palette
	v.generatePalette()

	return &v
}

func (v *FFTVisualizer) OnPropertyChanged(propertyName string) {
	switch propertyName {
	case "Luminosity":
		v.generatePalette()
	}
}

func (v *FFTVisualizer) generatePalette() {
	l := v.Luminosity
	for i := 0; i < FFT_PALETTE_SIZE; i++ {
		r, g, b := hueToRGB(1.0/3.0 - (float64(i)/768.0)/3.0)
		v.palette[i].Red = r * (float64(i) / FFT_PALETTE_SIZE) * l
		v.palette[i].Green = g * (float64(i) / FFT_PALETTE_SIZE) * l
		v.palette[i].Blue = b * (float64(i) / FFT_PALETTE_SIZE) * l
	}
}

func (v *FFTVisualizer) startPulseAudio() (io.ReadCloser, io.ReadCloser, error) {
	if v.cmd != nil {
		return nil, nil, fmt.Errorf("PulseAudio already started")
	}
	cmd := exec.Command("parec", "--device", v.pulseAudioDevice, "--format=s16le", "--channels=1", "--latency-msec=10")
	cmd.Env = v.pulseAudioOpts

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	err := cmd.Start()
	if err != nil {
		v.cmd = cmd
	}
	return stdout, stderr, err
}

func (v *FFTVisualizer) Start() {
	go v.Run()
}

func (v *FFTVisualizer) Run() {
	data, stderr, err := v.startPulseAudio()
	if err != nil {
		log.Printf("FFTVisualizer: error while startPulseAudio: %q", err.Error())
		return
	}

	// Compute Hann window
	hann := window.Hann(v.samplesPerFFT)

	// Allocate buffers
	vRead := make([]int16, v.samplesPerFrame)
	vFloat := make([]float64, v.samplesPerFFT)
	vFFT := make([]float64, v.samplesPerFFT)
	fMag := make([]float64, v.samplesPerFFT/2)
	leds := make([]Led, v.samplesPerFFT/2)

	for {
		// Read stdout as s16le values
		err = binary.Read(io.LimitReader(data, int64(2*v.samplesPerFrame)), binary.LittleEndian, &vRead)
		if err != nil {
			break
		}

		// Shift vFloat
		copy(vFloat, vFloat[v.samplesPerFrame:])

		// Append vRead at end of buffer as float
		for k, value := range vRead {
			vFloat[k+(v.samplesPerFFT-v.samplesPerFrame)] = float64(value)
		}

		// Copy vFloat to vFFT for FFT processing as hanning window is done in-place
		copy(vFFT, vFloat)

		// Apply a Hann window
		for i, w := range hann {
			vFFT[i] *= w
		}

		// Do FFT on those values
		f := fft.FFTReal(vFFT)

		// Convert it as magnitude
		for i := 0; i < v.samplesPerFFT/2; i++ {
			fMag[i] = cmplx.Abs(f[i])
		}

		// Generate graph
		for i := 0; i < v.samplesPerFFT/2; i++ {
			value := fMag[i] / 10000
			var color Led
			if value < 768 {
				color = v.palette[int(value)]
			} else {
				color = v.palette[FFT_PALETTE_SIZE-1]
			}
			leds[i] = color
		}

		v.SendData(leds)
	}

	// Print error if any
	if err != nil {
		log.Printf("FFTVisualizer: %q", err.Error())
	}

	// Print stderr if any
	e, _ := ioutil.ReadAll(stderr)
	if len(e) != 0 {
		log.Printf("FFTVisualizer: stderr: %q", string(e))
	}
}
