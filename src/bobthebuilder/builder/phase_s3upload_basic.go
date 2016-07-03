package builder
import (
  "github.com/mitchellh/goamz/aws"
  "github.com/mitchellh/goamz/s3"
  "bobthebuilder/logging"
  "bobthebuilder/config"
  "io/ioutil"
  "net/http"
  "errors"
  "time"
  "path"
  "os"
)

var AWS_CONFIG_ERR = errors.New("Invalid AWS configuration")
var S3_DEFAULT_ACL = s3.Private

type S3UploadPhase struct{
  BasicPhase
  Auth aws.Auth

  Region string
  Bucket string
  ACL s3.ACL
  FilenameToUpload string
  DestinationFileName string
}

func (p * S3UploadPhase)init(index int, wantACL string){
  p.Type = "S3UP_BASIC"
  p.StatusString = PHASE_STATUS_READY
  p.Index = index

  if p.DestinationFileName == ""{
    p.DestinationFileName = p.FilenameToUpload
  }
  if wantACL == "public"{
    p.ACL = s3.PublicRead
  } else {
    p.ACL = S3_DEFAULT_ACL
  }
}
func (p * S3UploadPhase)String()string{
  return "(phase)" + p.Type + " -- " + p.StatusString + " (" + p.FilenameToUpload + ")"
}






func (p* S3UploadPhase)phaseError(eCode int, statusString string)int{
  p.ErrorCode = eCode
  logging.Error("phase-s3-upload-basic", statusString)
  p.StatusString = statusString
  p.End = time.Now()
  p.Duration = p.End.Sub(p.Start)
  return eCode
}

func (p* S3UploadPhase)bucketExistsInBucketSet(bucketName string, buckets *s3.ListBucketsResp)bool{
  for _, bucket := range buckets.Buckets {
    if bucket.Name == bucketName{
      return true
    }
  }
  return false
}






func (p * S3UploadPhase)Run(r* Run, builder *Builder, defIndex int)int{
  var err error
  p.Start = time.Now()

  //run templates to sub in any variable information like dates etc
  p.DestinationFileName, err = ExecTemplate(p.DestinationFileName, p)
  if err != nil{
    p.WriteOutput( "Template Error (filename-destination): " + err.Error() + "\n", r, builder, defIndex)
    return p.phaseError(-11, "Template error")
  }
  p.FilenameToUpload, err = ExecTemplate(p.FilenameToUpload, p)
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

  //read file into a buffer
  pwd, _ := os.Getwd()
  p.WriteOutput( "Reading file: " + p.FilenameToUpload + "\n", r, builder, defIndex)
  data, err := ioutil.ReadFile(path.Join(pwd, BUILD_TEMP_FOLDER_NAME, p.FilenameToUpload))
  if err != nil {
    p.WriteOutput( "Unable to read source file: " + p.FilenameToUpload + "\n", r, builder, defIndex)
    p.WriteOutput( "Error: " + err.Error() + "\n", r, builder, defIndex)
    return p.phaseError(-4, "Could not read file")
  }

  //detect content type
  sizeToDetect := len(data)
  if sizeToDetect > 1024 {
    sizeToDetect = 1024
  }
  contType := http.DetectContentType(data[:sizeToDetect])


  p.WriteOutput( "Writing file of type " + contType + " to " + p.DestinationFileName + "...", r, builder, defIndex)
  err = s3cli.Bucket(p.Bucket).Put(p.DestinationFileName, data, contType, p.ACL)
  p.End = time.Now()
  p.Duration = p.End.Sub(p.Start)
  if err == nil{
    p.WriteOutput("DONE.\n", r, builder, defIndex)
    p.ErrorCode = 0
    p.StatusString = "Upload successful"
    return 0
  } else {
    p.WriteOutput("ERROR.\n", r, builder, defIndex)
    p.ErrorCode = -5
    p.StatusString = "Upload failed."
    return -5
  }
}
