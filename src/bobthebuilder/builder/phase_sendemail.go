package builder
import (
  "bobthebuilder/config"
  "bobthebuilder/logging"
  "html"
  "strconv"
  "strings"
  "time"
)

type SendEmailPhase struct{
  BasicPhase

  SendManual bool
  SendAllOutput bool
  Prefix string

  // special cases - control behaviour
  SendOnFailure bool
  SendOnSuccess bool

  Destinations []string
  SubjectOverride string
}

func (p * SendEmailPhase)init(index int){
  p.Type = "SEND_EMAIL"
  p.StatusString = PHASE_STATUS_READY
  p.Index = index
}

func (p * SendEmailPhase)String()string{
  return "(phase)" + p.Type + " -- " + p.StatusString
}

func (p* SendEmailPhase)phaseError(eCode int, statusString string)int{
  p.ErrorCode = eCode
  logging.Error("phase-send-email", statusString)
  p.StatusString = statusString
  p.End = time.Now()
  p.Duration = p.End.Sub(p.Start)
  return eCode
}



func (p * SendEmailPhase)Run(r* Run, builder *Builder, defIndex int)int{

  if !config.All().Gmail.Enable || config.All().Gmail.FromAddress == "" || config.All().Gmail.Password == "" {
    p.WriteOutput( "Gmail is not configured in the servers configuration. Is Enabled,\nFromAddress, and Password set?", r, builder, defIndex)
    return p.phaseError(-1, "Incomplete gmail configuration")
  }

  sent := false

  if p.SendOnFailure {
    if r.Status < STATUS_SUCCESS {
      p.WriteOutput( "Sending notification for failure condition.", r, builder, defIndex)
      errorCode, errorMsg := p.Send(r.Definition.Name + " failure", p.MakeLog("The definition failed to execute.", r, builder, defIndex, p.SendAllOutput), r, builder, defIndex)
      if errorCode != STATUS_SUCCESS {
        return p.phaseError(errorCode, errorMsg)
      }
      sent = true
    }
  }
  if p.SendOnSuccess {
    if r.Status == STATUS_SUCCESS {
      p.WriteOutput( "Sending notification for success condition.", r, builder, defIndex)
      errorCode, errorMsg := p.Send(r.Definition.Name + " success", p.MakeLog("The definition executed successfully.", r, builder, defIndex, p.SendAllOutput), r, builder, defIndex)
      if errorCode != STATUS_SUCCESS {
        return p.phaseError(errorCode, errorMsg)
      }
      sent = true
    }
  }
  if p.SendManual{
    var err error
    var prefix string
    prefix, err = ExecTemplate(p.Prefix, p, r, builder)
    if err != nil{
      p.WriteOutput( "Template Error (prefix): " + err.Error() + "\n", r, builder, defIndex)
    }

    content := prefix
    if p.SendAllOutput {
      content = p.MakeLog(prefix, r, builder, defIndex, p.SendAllOutput)
    } else {
      p.WriteOutput( "Content: " + content + "\n", r, builder, defIndex)
    }

    errorCode, errorMsg := p.Send(r.Definition.Name + " success", content, r, builder, defIndex)
    if errorCode != STATUS_SUCCESS {
      return p.phaseError(errorCode, errorMsg)
    }
    sent = true
  }

  p.End = time.Now()
  p.Duration = p.End.Sub(p.Start)
  p.ErrorCode = 0
  if !sent {
    p.WriteOutput( "No triggers - aborting send.", r, builder, defIndex)
    p.StatusString = "Send not required"
  } else {
    p.StatusString = "Send successful"
  }
  return 0
}

func (p * SendEmailPhase)Send(subject, content string, r* Run, builder *Builder, defIndex int)(errorcode int, errorMsg string){
  if p.SubjectOverride != "" {
    var err error
    subject, err = ExecTemplate(p.SubjectOverride, p, r, builder)
    if err != nil{
      p.WriteOutput( "Template Error (subject): " + err.Error() + "\n", r, builder, defIndex)
    }
    p.WriteOutput( "Subject: " + subject + "\n", r, builder, defIndex)
  }

  email := Compose(subject, content)
  email.From = config.All().Gmail.FromAddress
  email.Password = config.All().Gmail.Password
  if len(p.Destinations) == 0 && config.All().Gmail.DefaultToAddress != ""{
    email.AddRecipient(config.All().Gmail.DefaultToAddress)
  } else if len(p.Destinations) > 0{
    for _, dest := range p.Destinations {
      email.AddRecipient(dest)
    }
  } else {
     p.WriteOutput( "Error: No Destinations. Is config.Gmail.DefaultToAddress set?", r, builder, defIndex)
    return -4, "No destinations"
  }
  email.ContentType = "text/html; charset=utf-8"
  err := email.Send()
  if err != nil {
    p.WriteOutput( "Error: " + err.Error(), r, builder, defIndex)
    return -2, "Send Failure"
  }
  return STATUS_SUCCESS, ""
}


func (p * SendEmailPhase)MakeLog(prefix string, r* Run, builder *Builder, defIndex int, allOutput bool)string{
  var out string
  out += "<p>" + prefix + "</p>\n<p style=\"font-size: 21px;\"><b>Execution Summary</b></p>\n\n"
  for _, phase := range r.Phases {

    if phase.GetStatusString() == PHASE_STATUS_READY{
      continue
    }

    out += "<b style=\"font-size: 17px;\">" + statusStyled(phase.GetErrorCode(), html.EscapeString(phase.String())) + "</b>\n"
    out += "<table>"
    out += "<tr>"
    out += "<td style=\"font-weight: bold;\">Status</td><td><i>" + html.EscapeString(phase.GetStatusString()) + "</i></td>"
    out += "</tr>\n"

    t := phase.GetType()
    if t == "CLEAN" || t == "APT-CHECK" || t == "S3UP_BASIC" || t == "SET_ENV" || t == "TAR_TO_S3" || t == "BASE-INSTALL" || t == "SEND_EMAIL" || allOutput {
      out += "<tr>"
      out += "<td>&nbsp;</td><td>" + strings.Replace(html.EscapeString(strings.Join(phase.GetOutputs(), "<br>")), html.EscapeString("<br>"), "<br>", -1) + "</td>"
      out += "</tr>\n"
    }

    out += "<tr>"
    out += "<td style=\"font-weight: bold;\">Ex</td><td>" + emailErrorCodeStyled(phase.GetErrorCode()) + "</td>"
    out += "</tr>"
    out += "</table> <br><br>\n"
  }
  out += ""
  return out
}

func statusStyled(code int, status string)string{
  if code == 954321 {
    return "<span style=\"color: brown; text-decoration: line-through;\"> Execution skipped </span>"
  }
  if code == 0 {
    return status
  }
  return "<span style=\"color: red;\">" + status + "</span>"
}

func emailErrorCodeStyled(code int)string{
  if code == 954321 {
    return "<span style=\"color: green;\"> --- </span>"
  }

  if code == 0{
    return "<span style=\"color: green;\">0</span>"
  }
  return "<span style=\"color: red;\">" + strconv.Itoa(code) + "</span>"
}
