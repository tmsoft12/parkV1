package util

import "park/models/camera"

var validCameras = []camera.CameraType{
	camera.Inside,
	camera.Outside,
}

func IsValidCamera(camera camera.CameraType) bool {
	for _, validCamera := range validCameras {
		if camera == validCamera {
			return true
		}
	}
	return false
}
