package builder

import (
	"archive/tar"
	"bobthebuilder/config"
	"bobthebuilder/logging"
	"compress/gzip"
	"github.com/jhoonb/archivex"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/rlmcpherson/s3gof3r"
	"path"
	"strings"
	"time"
)

type TarToS3 struct {
	BasicPhase

	DestinationPath string
	Bucket          string
	Region          string
	Directories     []string
	Files           []string
}

func (p *TarToS3) init(index int) {
	p.Type = "TAR_TO_S3"
	p.StatusString = PHASE_STATUS_READY
	p.Index = index
}

func (p *TarToS3) String() string {
	return "(phase)" + p.Type + " -- " + p.StatusString
}

func (p *TarToS3) phaseError(eCode int, statusString string) int {
	p.ErrorCode = eCode
	logging.Error("phase-tar-to-s3", statusString)
	p.StatusString = statusString
	p.End = time.Now()
	p.Duration = p.End.Sub(p.Start)
	return eCode
}

func (p *TarToS3) bucketExistsInBucketSet(bucketName string, buckets *s3.ListBucketsResp) bool {
	for _, bucket := range buckets.Buckets {
		if bucket.Name == bucketName {
			return true
		}
	}
	return false
}

func (p *TarToS3) Run(r *Run, builder *Builder, defIndex int) int {
	var err error

	//run templates to sub in any variable information like dates etc
	p.DestinationPath, err = ExecTemplate(p.DestinationPath, p, r, builder)
	if err != nil {
		p.WriteOutput("Template Error (filename-destination): "+err.Error()+"\n", r, builder, defIndex)
		return p.phaseError(-11, "Template error")
	}

	//throw error if AWS not configured
	if !config.All().AWS.Enable || config.All().AWS.SecretKey == "" || config.All().AWS.AccessKey == "" {
		p.WriteOutput("Error: Invalid AWS configuration - have you enabled AWS in the config and provided the secret key / access key?\n", r, builder, defIndex)
		return p.phaseError(-1, AWS_CONFIG_ERR.Error())
	}

	var auth aws.Auth
	auth, err = aws.GetAuth(config.All().AWS.AccessKey, config.All().AWS.SecretKey)
	s3cli := s3.New(auth, aws.Regions[p.Region])

	//try and get bucket details - must have invalid creds if not -> error
	p.WriteOutput("Opening bucket: "+p.Bucket+" ("+p.Region+")\n", r, builder, defIndex)
	bList, err := s3cli.ListBuckets()
	if err != nil {
		p.WriteOutput("Error: "+err.Error()+"\n", r, builder, defIndex)
		return p.phaseError(-2, "Invalid AWS configuration")
	}

	//check the bucket exists
	if !p.bucketExistsInBucketSet(p.Bucket, bList) {
		p.WriteOutput("Error: Bucket does not exist in the specified region.\n", r, builder, defIndex)
		return p.phaseError(-3, "Bucket not found")
	}

	//lets setup the upload.
	keys := s3gof3r.Keys{AccessKey: config.All().AWS.AccessKey, SecretKey: config.All().AWS.SecretKey}
	s3Streamer := s3gof3r.New(strings.Replace(aws.Regions[p.Region].S3Endpoint, "https://", "", -1), keys)
	p.WriteOutput("S3 Endspoint: "+aws.Regions[p.Region].S3Endpoint+"\n", r, builder, defIndex)
	p.WriteOutput("Destination: "+p.DestinationPath+"\n", r, builder, defIndex)

	conf := &s3gof3r.Config{
		Concurrency: 2,
		PartSize:    6 * 1024 * 1024,
		NTry:        10,
		Md5Check:    true,
		Scheme:      "https",
		Client:      s3gof3r.ClientWithTimeout(20 * time.Second),
	}
	s3WriterObject, err := s3Streamer.Bucket(p.Bucket).PutWriter(p.DestinationPath, nil, conf)
	if err != nil {
		p.WriteOutput("Error: "+err.Error()+"\n", r, builder, defIndex)
		return p.phaseError(-4, "Transfer initialisation error")
	}

	gzipWriter := gzip.NewWriter(s3WriterObject)
	tF := archivex.TarFile{
		GzWriter:   gzipWriter,
		Writer:     tar.NewWriter(gzipWriter),
		Compressed: true,
		Name:       "<tar>",
	}

	for _, dir := range p.Directories {
		p.WriteOutput("Adding folder to archive: "+dir+"\n", r, builder, defIndex)
		err = tF.AddAll(path.Join(BuildDir, dir), true)
		if err != nil {
			p.WriteOutput("Error: "+err.Error()+"\n", r, builder, defIndex)
			tF.Close()
			s3WriterObject.Close()
			return p.phaseError(-10, "folder add error")
		}
	}

	for _, f := range p.Files {
		p.WriteOutput("Adding file to archive: "+f+"\n", r, builder, defIndex)
		err = tF.AddFile(path.Join(BuildDir, f))
		if err != nil {
			p.WriteOutput("Error: "+err.Error()+"\n", r, builder, defIndex)
			tF.Close()
			s3WriterObject.Close()
			return p.phaseError(-10, "file add error")
		}
	}

	err = tF.Close()
	if err != nil {
		p.WriteOutput("Error: "+err.Error()+"\n", r, builder, defIndex)
		s3WriterObject.Close()
		return p.phaseError(-5, "compression stream error")
	}

	err = s3WriterObject.Close()
	if err != nil {
		p.WriteOutput("Error: "+err.Error()+"\n", r, builder, defIndex)
		return p.phaseError(-6, "s3 stream error")
	}

	p.End = time.Now()
	p.Duration = p.End.Sub(p.Start)
	p.ErrorCode = 0
	p.StatusString = "Upload successful"
	return 0
}
