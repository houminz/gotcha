module github.com/SimpCosm/crddemo

go 1.14

require (
	k8s.io/api v0.0.0-20201020200614-54bcc311e327
	k8s.io/apimachinery v0.0.0-20201020200440-554eef9dbf66
	k8s.io/client-go v0.0.0-20201020200834-d1a4fe5f2d96
	k8s.io/code-generator v0.0.0-20201020200306-60862b8acf58
	k8s.io/klog/v2 v2.2.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20201020200614-54bcc311e327
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20201020200440-554eef9dbf66
	k8s.io/client-go => k8s.io/client-go v0.0.0-20201020200834-d1a4fe5f2d96
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20201020200306-60862b8acf58
)
