package main

import "fmt"

func main() {
	var l []string

	if l == nil {
		fmt.Println("is nil")
	}

}

//
//func main() {
//	session, err := aws_session.NewSession()
//	if err != nil {
//		panic(err)
//	}
//	s3 := aws_s3.New(session)
//
//	out, err := s3.ListObjectsV2(&aws_s3.ListObjectsV2Input{
//		S3BlobStorage:              aws.String("serverless-monorepo-101-deployments"),
//		ContinuationToken:   nil,
//		Delimiter:           nil,
//		EncodingType:        nil,
//		ExpectedBucketOwner: nil,
//		FetchOwner:          nil,
//		MaxKeys:             nil,
//		//Prefix:              aws.String("demo-service/"),
//		RequestPayer: nil,
//		StartAfter:   nil,
//	})
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(*out.KeyCount)
//
//	//for _, o := range out.Contents {
//	//	fmt.Println(*o.Key)
//	//	fmt.Println(o.LastModified)
//	//	fmt.Println("....")
//	//}
//	//
//	//fmt.Println(lastestObject(out.Contents))
//}
//
//func lastestObject(objects []*aws_s3.Object) *aws_s3.Object {
//	var latest *aws_s3.Object
//
//	for _, obj := range objects {
//		if latest == nil {
//			latest = obj
//			fmt.Println("nil updated!")
//			continue
//		}
//		if obj.LastModified.After(*latest.LastModified) {
//			latest = obj
//			fmt.Println("compare updated!")
//			continue
//		}
//	}
//	return latest
//}
