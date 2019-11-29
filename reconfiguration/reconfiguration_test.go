package reconfiguration

import (
	"testing"
	"time"
)

func TestReconfiguration_GenerateReconfiguration(t *testing.T) {
	re:=NewReconfiguration()
	time1:=time.Now()
	time.Sleep(time.Second*2)
	re.GenerateReconfiguration(time1)
}
func TestReconfiguration_GenerateCredibleShardPlan(t *testing.T) {
	re:=NewReconfiguration()
	re.GenerateCredibleShardPlan()
}
func TestReconfiguration_GenerateCrediblePlan(t *testing.T) {
	re:=NewReconfiguration()
	re.GenerateCrediblePlan()
}
func TestReconfiguration_GetNodeInfo(t *testing.T) {
	re:=NewReconfiguration()
	re.GetNodeInfo()
}
func TestReconfiguration_SendReconfiguration(t *testing.T) {
	re:=NewReconfiguration()
	re.SendReconfiguration()
}
