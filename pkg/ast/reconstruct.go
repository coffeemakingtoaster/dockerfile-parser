package ast

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func formatIfValue(fstring, value string) string {
	if len(value) == 0 {
		return ""
	}
	return fmt.Sprintf(fstring, value)
}

func escapeSlice[T any](slice []T) string {
	jsonBytes, err := json.Marshal(slice)
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

func (sn *StageNode) Reconstruct() []string {
	reconstructed := []string{}
	if sn.Image != "" {
		fromInstruction := fmt.Sprintf("FROM %s", sn.Image)
		if sn.Name != "" {
			fromInstruction += fmt.Sprintf(" AS %s", sn.Name)
		}
		reconstructed = append(reconstructed, fromInstruction)
	}
	for _, instructionNode := range sn.Instructions {
		reconstructed = append(reconstructed, instructionNode.Reconstruct()...)
	}
	if sn.Subsequent == nil {
		return reconstructed
	}
	return append(reconstructed, sn.Subsequent.Reconstruct()...)
}

func (ai *AddInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s ", ai.Instruction())
	reconstructed += formatIfValue("--keep-git-dir=%s ", strconv.FormatBool(ai.KeepGitDir))
	reconstructed += formatIfValue("--checksum=%s ", ai.CheckSum)
	reconstructed += formatIfValue("--chown=%s ", ai.Chown)
	reconstructed += formatIfValue("--chmod=%s ", ai.Chmod)
	reconstructed += formatIfValue("--link=%s ", strconv.FormatBool(ai.Link))
	reconstructed += formatIfValue("--exclude=%s ", ai.Exclude)
	reconstructed += fmt.Sprintf("%s ", strings.Join(ai.Source, " "))
	reconstructed += fmt.Sprintf("%s", ai.Destination)
	return []string{reconstructed}
}

func (ai *ArgInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s %s=%s", ai.Instruction(), ai.Name, ai.Value)
	return []string{reconstructed}
}

func (ci *CmdInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s %s", ci.Instruction(), escapeSlice(ci.Cmd))
	return []string{reconstructed}
}
func (ci *CopyInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s ", ci.Instruction())

	reconstructed += formatIfValue("--keep-git-dir=%s ", strconv.FormatBool(ci.KeepGitDir))
	reconstructed += formatIfValue("--chown=%s ", ci.Chown)
	reconstructed += formatIfValue("--link=%s ", strconv.FormatBool(ci.Link))
	reconstructed += formatIfValue("--from=%s ", ci.From)
	reconstructed += fmt.Sprintf("%s ", strings.Join(ci.Source, " "))
	reconstructed += fmt.Sprintf("%s", ci.Destination)
	return []string{reconstructed}
}
func (ei *EntrypointInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s %s", ei.Instruction(), escapeSlice(ei.Exec))
	return []string{reconstructed}
}
func (ei *EnvInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s", ei.Instruction())
	for k, v := range ei.Pairs {
		reconstructed += fmt.Sprintf(" %s=%s", k, v)
	}
	return []string{reconstructed}
}

func (ei *ExposeInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s", ei.Instruction())
	for _, port := range ei.Ports {
		protocol := "tcp"
		if !port.IsTCP {
			protocol = "udp"
		}
		reconstructed += fmt.Sprintf(" %s/%s", port.Port, protocol)
	}
	return []string{reconstructed}
}
func (hi *HealthcheckInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s ", hi.Instruction())
	if hi.CancelStatement {
		return []string{reconstructed + "NONE"}
	}
	reconstructed += formatIfValue("--interval=%s ", hi.Interval)
	reconstructed += formatIfValue("--timeout=%s ", hi.Timeout)
	reconstructed += formatIfValue("--start-period=%s ", hi.StartPeriod)
	reconstructed += formatIfValue("--start-interval=%s ", hi.StartInterval)
	reconstructed += formatIfValue("--retries=%s ", strconv.Itoa(hi.Retries))
	reconstructed += escapeSlice(hi.Cmd)
	return []string{reconstructed}
}
func (li *LabelInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s", li.Instruction())
	for k, v := range li.Pairs {
		reconstructed += fmt.Sprintf(" %s=%s", k, v)
	}
	return []string{reconstructed}
}
func (mi *MaintainerInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s %s", mi.Instruction(), mi.Name)
	return []string{reconstructed}
}
func (oi *OnbuildInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s ", oi.Instruction())
	nested := oi.Trigger.Reconstruct()
	nested[0] = reconstructed + nested[0]
	return nested
}
func (ri *RunInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s ", ri.Instruction())
	if !ri.ShellForm && !ri.IsHeredoc {
		return []string{reconstructed + escapeSlice(strings.Split(ri.Cmd[0], " "))}
	}
	if ri.IsHeredoc {
		reconstructed += "<< "
	}
	ri.Cmd[0] = reconstructed + ri.Cmd[0]
	return ri.Cmd
}
func (si *ShellInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s %s", si.Instruction(), escapeSlice(si.Shell))
	return []string{reconstructed}
}

func (si *StopsignalInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s %s", si.Instruction(), si.Signal)
	return []string{reconstructed}
}

func (ui *UserInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s %s", ui.Instruction(), ui.User)
	return []string{reconstructed}
}

func (vi *VolumeInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s %s", vi.Instruction(), escapeSlice(vi.Mounts))
	return []string{reconstructed}
}

func (wi *WorkdirInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s %s", wi.Instruction(), wi.Path)
	return []string{reconstructed}
}
func (ci *CommentInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("# %s", ci.Text)
	return []string{reconstructed}
}

// For the edge case that instruction supplied to ONBUILD cannot be parsed
func (ui *UnknownInstructionNode) Reconstruct() []string {
	return []string{ui.Text}
}
