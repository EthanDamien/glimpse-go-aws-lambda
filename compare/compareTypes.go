package compare

type CompareResponse struct {
	DESC string `json:"body"`
}

// https://docs.aws.amazon.com/rekognition/latest/APIReference/API_CompareFaces.html
const CompareRequestSyntaxTest = `{
	"QualityFilter": "HIGH",
	"SimilarityThreshold": 95,
	"SourceImage": { 
	   "Bytes": blob,
	   "S3Object": { 
		  "Bucket": "facefiles",
		  "Name": "%s",
	   }
	},
	"TargetImage": { 
	   "Bytes": blob,
	   "S3Object": { 
		  "Bucket": "facefiles",
		  "Name": "%s",
	   }
	}
 }
 `
