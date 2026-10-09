time=2026-10-09T11:52:02.085+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
--- FAIL: TestGoldenSingleToolCall (0.01s)
    loop_golden_test.go:81: call identity = {TaskID:fbb95865-367d-4b21-ae9c-b3e4c7e2a539 CorrelationID:fbb95865-367d-4b21-ae9c-b3e4c7e2a539 CallID:call_e1 Name:echo Args:{"text":"22 摄氏度，晴"} Timeout:2s}, want task fbb95865-367d-4b21-ae9c-b3e4c7e2a539 and correlation id fbb95865-367d-4b21-ae9c-b3e4c7e2a539#call_e1 (task id prefix + this call's own id, never the bare task id)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/agent	0.040s
FAIL
MUTANT_TEST_RC=1
