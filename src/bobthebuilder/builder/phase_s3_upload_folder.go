package builder
import (
  "github.com/mitchellh/goamz/aws"
  "github.com/mitchellh/goamz/s3"
  "bobthebuilder/logging"
  "bobthebuilder/config"
  "io/ioutil"
  "net/http"
  "time"
  "path"
  "fmt"
  "os"
)

type S3UploadFolderPhase struct{
  BasicPhase
  Auth aws.Auth

  Region string
  Bucket string
  ACL s3.ACL
  SourceFolder string
  DestinationFolder string
}

func (p * S3UploadFolderPhase)init(index int, wantACL string){
  p.Type = "S3_UPLOAD_FOLDER"
  p.StatusString = PHASE_STATUS_READY
  p.Index = index

  if wantACL == "public"{
    p.ACL = s3.PublicRead
  } else {
    p.ACL = S3_DEFAULT_ACL
  }
}
func (p * S3UploadFolderPhase)String()string{
  return "(phase)" + p.Type + " -- " + p.StatusString + " (" + p.SourceFolder + ")"
}






func (p* S3UploadFolderPhase)phaseError(eCode int, statusString string)int{
  p.ErrorCode = eCode
  logging.Error("phase-s3-upload-folder", statusString)
  p.StatusString = statusString
  p.End = time.Now()
  p.Duration = p.End.Sub(p.Start)
  return eCode
}

func (p* S3UploadFolderPhase)bucketExistsInBucketSet(bucketName string, buckets *s3.ListBucketsResp)bool{
  for _, bucket := range buckets.Buckets {
    if bucket.Name == bucketName{
      return true
    }
  }
  return false
}






func (p * S3UploadFolderPhase)Run(r* Run, builder *Builder, defIndex int)int{
  var err error

  //run templates to sub in any variable information like dates etc
  p.DestinationFolder, err = ExecTemplate(p.DestinationFolder, p, r, builder)
  if err != nil{
    p.WriteOutput( "Template Error (filename-destination): " + err.Error() + "\n", r, builder, defIndex)
    return p.phaseError(-11, "Template error")
  }
  p.SourceFolder, err = ExecTemplate(p.SourceFolder, p, r, builder)
  if err != nil{
    p.WriteOutput( "Template Error (filename): " + err.Error() + "\n", r, builder, defIndex)
    return p.phaseError(-11, "Template error")
  }

  //throw error if AWS not configured
  if !config.All().AWS.Enable || config.All().AWS.SecretKey == "" || config.All().AWS.AccessKey == "" {
    p.WriteOutput( "Error: Invalid AWS configuration - have you enabled AWS in the config and provided the secret key / access key?\n", r, builder, defIndex)
    return p.phaseError(-1, AWS_CONFIG_ERR.Error())
  }

  p.Auth, err = aws.GetAuth(config.All().AWS.AccessKey, config.All().AWS.SecretKey)
  s3cli := s3.New(p.Auth, aws.Regions[p.Region])

  //try and get bucket details - must have invalid creds if not -> error
  p.WriteOutput("Opening bucket: " + p.Bucket + " (" + p.Region + ")\n", r, builder, defIndex)
  bList, err := s3cli.ListBuckets()
  if err != nil {
    p.WriteOutput( "Error: " + err.Error() + "\n", r, builder, defIndex)
    return p.phaseError(-2, "Invalid AWS configuration")
  }

  //check the bucket exists
  if !p.bucketExistsInBucketSet(p.Bucket, bList) {
    p.WriteOutput( "Error: Bucket does not exist in the specified region.\n", r, builder, defIndex)
    return p.phaseError(-3, "Bucket not found")
  }

  pwd, _ := os.Getwd()
  files, err := ioutil.ReadDir(path.Join(pwd, BUILD_TEMP_FOLDER_NAME, p.SourceFolder))
  if err != nil {
    p.WriteOutput( "Error: " + err.Error() + "\n", r, builder, defIndex)
    return p.phaseError(-6, "Internal error")
  }
  for _, file := range files {
    if file.IsDir() {
      continue
    }
    p.WriteOutput( "Reading file: " + file.Name() + fmt.Sprintf(" (%.4f Mb)", float64(file.Size())/1024.0/1024.0) + "\n", r, builder, defIndex)
    data, err := ioutil.ReadFile(path.Join(pwd, BUILD_TEMP_FOLDER_NAME, file.Name()))
    if err != nil {
      p.WriteOutput( "Unable to read source file: " + file.Name() + "\n", r, builder, defIndex)
      p.WriteOutput( "Error: " + err.Error() + "\n", r, builder, defIndex)
      return p.phaseError(-4, "Could not read file")
    }

    //detect content type
    sizeToDetect := len(data)
    if sizeToDetect > 1024 {
      sizeToDetect = 1024
    }
    contType := http.DetectContentType(data[:sizeToDetect])


    p.WriteOutput( "Writing file of type " + contType + " to " + path.Join(p.DestinationFolder, file.Name()) + "...", r, builder, defIndex)
    err = s3cli.Bucket(p.Bucket).Put(path.Join(p.DestinationFolder, file.Name()), data, contType, p.ACL)
    if err == nil{
      p.WriteOutput("DONE.\n", r, builder, defIndex)
    } else {
      p.WriteOutput("ERROR.\n", r, builder, defIndex)
      p.WriteOutput("Error: " + err.Error() + "\n", r, builder, defIndex)
      p.End = time.Now()
      p.Duration = p.End.Sub(p.Start)
      p.ErrorCode = -5
      p.StatusString = "Upload failed."
      return -5
    }
  }
  p.End = time.Now()
  p.Duration = p.End.Sub(p.Start)
  p.ErrorCode = 0
  p.StatusString = "Upload successful"
  return 0
}
