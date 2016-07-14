package util

import (
  "bobthebuilder/logging"
  "os/exec"
  "strings"
  "errors"
)

//handles data fetch and some of the presentation for remote-dropdowns.



func GetGitLookupData(options map[string]interface{})([]map[string]interface{}, error){
    gitURL, ok := options["git-url"].(string)
    if !ok{
      return nil, errors.New("Specify a git-url in the parameters options.")
    }
    branchNamesOnly, _ := options["branchNamesOnly"].(bool)

    out, err := exec.Command("git", "ls-remote", "--heads", gitURL).Output()
    if err != nil {
      logging.Error("web-definitions-api", err)
      return nil, err
    }

    branches := []map[string]interface{}{}
    lines := strings.Split(string(out), "\n")
    for _, line := range lines{
      if len(line) <= 4{
        continue
      }
      spl := strings.Split(line, "\t")
      spl[1] = strings.TrimPrefix(spl[1], "refs/heads/")
      val := spl[0]
      if branchNamesOnly{
        val = spl[1]
      }
      branches = append(branches, map[string]interface{}{"name": spl[1], "value": val,})
    }
    return branches, nil
}
