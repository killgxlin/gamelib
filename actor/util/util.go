package util

import (
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func KeyToPID(key string) *actor.PID {
	pid := actor.NewLocalPID("")
	i := strings.Index(key, "#")
	pid.Id = key[i+1:]
	if i > 0 {
		pid.Address = key[0:i]
	}
	return pid
}

func PIDToKey(pid *actor.PID) string {
	return pid.Address + "#" + pid.Id
}

func KeyToPID2(key string) *actor.PID {
	if key == "" {
		return nil
	}
	pid := actor.NewLocalPID("")
	i := strings.Index(key, "#")
	pid.Id = key[i+1:]
	if i > 0 {
		pid.Address = key[0:i]
	}
	return pid
}

func PIDToKey2(pid *actor.PID) string {
	if pid == nil {
		return ""
	}
	return pid.Address + "#" + pid.Id
}
