package effect

import (
	"fmt"
	"ledfx/color"
	"ledfx/config"
	"ledfx/device"
	"ledfx/logger"
	"math"
	"time"
)

type Effect interface {
	AssembleFrame(phase float64, ledCount int, effectColor color.Color) (colors []color.Color)
}

type EffectConfig struct {
	Blur       float64
	Flip       bool
	Mirror     bool
	Brightness float32
	Background color.Color
}

func StartEffect(deviceConfig config.DeviceConfig, effect Effect, fps int, done <-chan bool) error {
	logger.Logger.Debug(fmt.Sprintf("fps: %v", fps))
	usPerFrame := (float64(1.0) / float64(fps))
	usPerFrameDuration := time.Duration(usPerFrame*1000000.0) * time.Microsecond
	logger.Logger.Debug(fmt.Sprintf("usPerFrameDuration: %v", usPerFrameDuration.Microseconds()))
	ticker := time.NewTicker(usPerFrameDuration)
	phase := 0.0 // phase of the effect (range 0.0 to 2π)

	// TODO: choose type of device dynamically based on the deviceConfig
	var device = &device.UdpDevice{
		Name:     deviceConfig.Name,
		Port:     deviceConfig.Port,
		Protocol: device.UdpProtocols[deviceConfig.UdpPacketType],
		Config:   deviceConfig,
	}

	err := device.Init()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	defer ticker.Stop()

	// TODO: this should be in effect config
	speed := 1.0 // beats per minute

	for {
		select {
		case <-done:
			fmt.Println("Done!")
			device.Close()
			return nil
		case <-ticker.C:
			// TODO: get pixelCount and color from config
			// TODO: this should be
			newColor, err := color.NewColor(color.LedFxColors["red"])
			if err != nil {
				return err
			}
			logger.Logger.Debug(fmt.Sprintf("phase: %v", phase))
			err = device.SendData(effect.AssembleFrame(phase, device.Config.PixelCount, newColor), 0xff)
			if err != nil {
				return err
			}
			// Increment the phase (range: 0 - 2π)
			phase += ((2 * math.Pi) / float64(fps)) * speed
			if phase >= (2 * math.Pi) {
				phase = 0.0
			}
		}
	}
}

// TODO: StopEffect
