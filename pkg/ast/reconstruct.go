package ast

import (
	"encoding/json"
	"fmt"
	"slices"
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
		var fromInstruction strings.Builder
		fromInstruction.WriteString(fmt.Sprintf("FROM %s", sn.Image))
		if sn.Name != "" {
			fromInstruction.WriteString(fmt.Sprintf(" AS %s", sn.Name))
		}
		reconstructed = append(reconstructed, fromInstruction.String())
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
	reconstructed := fmt.Sprintf("%s", ai.Instruction())
	keys := make([]string, len(ai.Pairs))
	index := 0
	for k := range ai.Pairs {
		keys[index] = k
		index++
	}
	slices.Sort(keys)
	for _, k := range keys {
		reconstructed += fmt.Sprintf(" %s=%s", k, ai.Pairs[k])
	}
	return []string{reconstructed}
}

func (ci *CmdInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s %s", ci.Instruction(), escapeSlice(ci.Cmd))
	return []string{reconstructed}
}
func (ci *CopyInstructionNode) Reconstruct() []string {
	var reconstructed strings.Builder
	reconstructed.WriteString(fmt.Sprintf("%s ", ci.Instruction()))

	reconstructed.WriteString(formatIfValue("--keep-git-dir=%s ", strconv.FormatBool(ci.KeepGitDir)))
	reconstructed.WriteString(formatIfValue("--chown=%s ", ci.Chown))
	reconstructed.WriteString(formatIfValue("--link=%s ", strconv.FormatBool(ci.Link)))
	reconstructed.WriteString(formatIfValue("--from=%s ", ci.From))
	reconstructed.WriteString(fmt.Sprintf("%s ", strings.Join(ci.Source, " ")))
	reconstructed.WriteString(fmt.Sprintf("%s", ci.Destination))
	return []string{reconstructed.String()}
}
func (ei *EntrypointInstructionNode) Reconstruct() []string {
	reconstructed := fmt.Sprintf("%s %s", ei.Instruction(), escapeSlice(ei.Exec))
	return []string{reconstructed}
}
func (ei *EnvInstructionNode) Reconstruct() []string {
	var reconstructed strings.Builder
	reconstructed.WriteString(fmt.Sprintf("%s", ei.Instruction()))
	keys := make([]string, len(ei.Pairs))
	index := 0
	for k := range ei.Pairs {
		keys[index] = k
		index++
	}
	slices.Sort(keys)
	for _, k := range keys {
		reconstructed.WriteString(fmt.Sprintf(" %s=%s", k, ei.Pairs[k]))
	}
	return []string{reconstructed.String()}
}

func (ei *ExposeInstructionNode) Reconstruct() []string {
	var reconstructed strings.Builder
	reconstructed.WriteString(fmt.Sprintf("%s", ei.Instruction()))
	for _, port := range ei.Ports {
		protocol := "tcp"
		if !port.IsTCP {
			protocol = "udp"
		}
		reconstructed.WriteString(fmt.Sprintf(" %s/%s", port.Port, protocol))
	}
	return []string{reconstructed.String()}
}
func (hi *HealthcheckInstructionNode) Reconstruct() []string {
	var reconstructed strings.Builder
	reconstructed.WriteString(fmt.Sprintf("%s ", hi.Instruction()))
	if hi.CancelStatement {
		reconstructed.WriteString("NONE")
		return []string{reconstructed.String()}
	}
	reconstructed.WriteString(formatIfValue("--interval=%s ", hi.Interval))
	reconstructed.WriteString(formatIfValue("--timeout=%s ", hi.Timeout))
	reconstructed.WriteString(formatIfValue("--start-period=%s ", hi.StartPeriod))
	reconstructed.WriteString(formatIfValue("--start-interval=%s ", hi.StartInterval))
	reconstructed.WriteString(formatIfValue("--retries=%s ", strconv.Itoa(hi.Retries)))
	reconstructed.WriteString(escapeSlice(hi.Cmd))
	return []string{reconstructed.String()}
}
func (li *LabelInstructionNode) Reconstruct() []string {
	var reconstructed strings.Builder
	reconstructed.WriteString(fmt.Sprintf("%s", li.Instruction()))
	keys := make([]string, len(li.Pairs))
	index := 0
	for k := range li.Pairs {
		keys[index] = k
		index++
	}
	slices.Sort(keys)
	for _, k := range keys {
		reconstructed.WriteString(fmt.Sprintf(" %s=%s", k, li.Pairs[k]))
	}
	return []string{reconstructed.String()}
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
	var reconstructed strings.Builder
	reconstructed.WriteString(fmt.Sprintf("%s ", ri.Instruction()))
	if !ri.ShellForm && !ri.IsHeredoc {
		reconstructed.WriteString(escapeSlice(ri.Cmd))
		return []string{reconstructed.String()}
	}
	if ri.IsHeredoc {
		reconstructed.WriteString("<< ")
	}
	ri.Cmd[0] = reconstructed.String() + ri.Cmd[0]
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

func (*EmptyLineNode) Reconstruct() []string {
	return []string{""}
}
